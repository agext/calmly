// Copyright 2015 ALRUX Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package calmly

import (
	"fmt"
	"strings"
	"testing"
)

type mockLogger struct {
	log string
}

func (ml *mockLogger) Print(s ...interface{}) {
	ml.log += fmt.Sprintln(s...)
}
func (ml *mockLogger) Fatal(s ...interface{}) {
	ml.log += "[FATAL] " + fmt.Sprintln(s...)
}
func (ml *mockLogger) Panic(s ...interface{}) {
	ml.log += "[PANIC] " + fmt.Sprintln(s...)
}

func TestLevelNames(t *testing.T) {
	for level, name := range map[int8]string{
		OK:    "OK",
		ERROR: "ERROR",
		PANIC: "PANIC",
		FATAL: "FATAL",
		17:    "?",
	} {
		if levelName(level) != name {
			t.Errorf(`levelName(%d) = %q, want %q`, level, levelName(level), name)
		}
	}
}

func TestSetters(t *testing.T) {
	out := &Outcome{}
	if ol := out.Level(); ol != OK {
		t.Errorf(`default.Level() = %q (%d), want %q`, levelName(ol), ol, levelName(OK))
	}
	if out.Error() != "" {
		t.Errorf(`default.Error() = %q, want %q`, out.Error(), "")
	}
	if ol := out.SetLevel(FATAL).Level(); ol != FATAL {
		t.Errorf(`SetLevel(FATAL).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(FATAL))
	}
	if ol := out.SetLevel(17).Level(); ol != FATAL {
		t.Errorf(`SetLevel(17).Level() = %q (%d), want %q (unchanged previous value)`, levelName(ol), ol, levelName(FATAL))
	}
	if out.SetCode(17).Code() != 17 {
		t.Errorf(`SetCode(17).Code() = 0x%04x, want 0x%04x`, out.Code(), 17)
	}
	if out.SetText("xyz").Text() != "xyz" {
		t.Errorf(`SetText("xyz").Text() = %q, want %q`, out.Text(), "xyz")
	}
	info := out.AddInfo("line 1", "line 2", "debug.stack").Info()
	if len(info) != 3 {
		t.Errorf(`len(AddInfo("line 1", "line 2", "debug.stack").Info()) = %q, want %q`, len(info), 3)
	} else {
		if info[0] != "line 1" {
			t.Errorf(`AddInfo("line 1", "line 2", "debug.stack").Info()[0] = %q, want %q`, info[0], "line 1")
		}
		if info[1] != "line 2" {
			t.Errorf(`AddInfo("line 1", "line 2", "debug.stack").Info()[1] = %q, want %q`, info[1], "line 2")
		}
		if !strings.Contains(info[2], "goroutine") || !strings.Contains(info[2], "calmly.TestSetters") {
			t.Errorf(`AddInfo("line 1", "line 2", "debug.stack").Info()[2] does not contain stack trace (got %q)`, info[2])
		}
	}
}

func TestLog(t *testing.T) {
	log := &mockLogger{}
	out := &Outcome{val: 17, err: fmt.Errorf("test"), text: "abc"}
	out.Log(log).SetLevel(ERROR).Log(log).SetLevel(PANIC).Log(log).SetLevel(FATAL).SetCode(17).Log(log)
	if log.log != "abc\n[PANIC] abc\n[FATAL] abc (code: 0x0011)\n" {
		t.Errorf(`logging test got %q, want %q`, log.log, "abc\n[PANIC] abc\n[FATAL] abc (code: 0x0011)\n")
	}
}

func TestStack(t *testing.T) {
	assertTryPanic := func(out *Outcome, action, text string) {
		oc := out.Code()
		if oc != ERR_TRY_PANIC {
			t.Errorf(action+`.Code() = 0x%04x, want 0x%04x`, oc, ERR_TRY_PANIC)
		}
		ot := out.Text()
		if !strings.HasPrefix(ot, "panic: ") {
			t.Errorf(action+`.Text() does not begin with "panic: " (got %q)`, ot)
		}
		if !strings.Contains(ot, text) {
			t.Errorf(action+`.Text() does not contain %q (got %q)`, text, ot)
		}
		ov := out.Value()
		if ov != nil {
			t.Errorf(action+`.Value() = %v, want %v`, ov, nil)
		}
		oe := out.Err()
		if oe != nil {
			t.Errorf(action+`.Err() = %v, want %v`, oe, nil)
		}
		if orv, ore := out.Result(); orv != ov || ore != oe {
			t.Errorf(action+`.Result() should equal (`+action+`.Value(), `+action+`.Err()); got (%v, %v != %v, $v)`, orv, ore, ov, oe)
		}
		if oes, exp := out.Error(), ot+fmt.Sprintf(` (code: 0x%04x)`, oc); oes != exp {
			t.Errorf(action+`.Error() = %q, want %q`, oes, exp)
		}
		info := out.info
		if len(info) != 1 {
			t.Errorf(`len(`+action+`.Info()) = %q, want %q`, len(info), 1)
		} else {
			if !strings.Contains(info[0], "goroutine") || !strings.Contains(info[0], "calmly.TestStack") {
				t.Errorf(action+`.Info()[0] does not contain stack trace (got %q)`, info[0])
			}
		}
	}
	divByZero := func() {
		a, b := 1, 0
		_ = a / b
	}

	out := Try(divByZero)
	if ol := out.Level(); ol != PANIC {
		t.Errorf(`Try(divByZero).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(PANIC))
	}
	assertTryPanic(out, `Try(divByZero)`, `divide by zero`)
	caught := false
	out.Catch(func(o *Outcome) {
		caught = true
		assertTryPanic(o, `Try(divByZero/*inside Catch()*/)`, `divide by zero`)
	})
	if !caught {
		t.Errorf(`Try(divByZero).Catch(f) should call f(*Outcome) on PANIC`)
	}
	if ol := out.Level(); ol != PANIC {
		t.Errorf(`Try(divByZero).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(PANIC))
	}
	assertTryPanic(out, `Try(divByZero).Catch()`, `divide by zero`)
	out.KeepCalm()
	if ol := out.Level(); ol != ERROR {
		t.Errorf(`Try(divByZero).KeepCalm().Level() = %q (%d), want %q`, levelName(ol), ol, levelName(ERROR))
	}
	assertTryPanic(out, `Try(divByZero).KeepCalm()`, `divide by zero`)

	out = Try(divByZero).Escalate()
	if ol := out.Level(); ol != FATAL {
		t.Errorf(`Try(divByZero).Escalate().Level() = %q (%d), want %q`, levelName(ol), ol, levelName(FATAL))
	}
	assertTryPanic(out, `Try(divByZero).Escalate()`, `divide by zero`)

	out = Try(func() error {
		divByZero()
		return fmt.Errorf("divByZero should panic")
	})
	if ol := out.Level(); ol != PANIC {
		t.Errorf(`Try(divByZeroErr).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(PANIC))
	}
	assertTryPanic(out, `Try(divByZeroErr)`, `divide by zero`)

	out = Try(func() interface{} {
		divByZero()
		return 17
	})
	if ol := out.Level(); ol != PANIC {
		t.Errorf(`Try(divByZeroVal).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(PANIC))
	}
	assertTryPanic(out, `Try(divByZeroVal)`, `divide by zero`)

	out = Try(func() (interface{}, error) {
		divByZero()
		return 17, fmt.Errorf("divByZero should panic")
	})
	if ol := out.Level(); ol != PANIC {
		t.Errorf(`Try(divByZeroValErr).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(PANIC))
	}
	assertTryPanic(out, `Try(divByZeroValErr)`, `divide by zero`)

	out = Try(func() (interface{}, error) {
		return 17, nil
	})
	if ol := out.Level(); ol != OK {
		t.Errorf(`Try(goodFunc).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(OK))
	}
	oc := out.Code()
	if oc != 0 {
		t.Errorf(`Try(goodFunc).Code() = 0x%04x, want 0x%04x`, oc, 0)
	}
	ot := out.Text()
	if ot != "" {
		t.Errorf(`Try(goodFunc).Text() = %q, want %q`, ot, "")
	}
	ov := out.Value()
	if ov.(int) != 17 {
		t.Errorf(`Try(goodFunc).Value() = %v, want %v`, ov, 17)
	}
	oe := out.Err()
	if oe != nil {
		t.Errorf(`Try(goodFunc).Err() = %v, want %v`, oe, nil)
	}
	if orv, ore := out.Result(); orv != ov || ore != oe {
		t.Errorf(`Try(goodFunc).Result() should equal (Try(goodFunc).Value(), Try(goodFunc).Err()); got (%v, %v != %v, %v)`, orv, ore, ov, oe)
	}
	if oes, exp := out.Error(), ""; oes != exp {
		t.Errorf(`Try(goodFunc).Error() = %q, want %q`, oes, exp)
	}
	info := out.info
	if len(info) != 0 {
		t.Errorf(`len(Try(goodFunc).Info()) = %q, want %q`, len(info), 0)
	}

	out = Try(func() (int, error) {
		return 17, nil
	})
	if ol := out.Level(); ol != ERROR {
		t.Errorf(`Try(badFunc).Level() = %q (%d), want %q`, levelName(ol), ol, levelName(ERROR))
	}
	oc = out.Code()
	if oc != ERR_TRY_ARG {
		t.Errorf(`Try(badFunc).Code() = 0x%04x, want 0x%04x`, oc, ERR_TRY_ARG)
	}
	ot = out.Text()
	if !strings.HasPrefix(ot, "Try: unsupported argument type") {
		t.Errorf(`Try(badFunc).Text() does not begin with "Try: unsupported argument type" (got %q)`, ot)
	}
	ov = out.Value()
	if ov != nil {
		t.Errorf(`Try(badFunc).Value() = %v, want %v`, ov, nil)
	}
	oe = out.Err()
	if oe != nil {
		t.Errorf(`Try(badFunc).Err() = %v, want %v`, oe, nil)
	}
	if orv, ore := out.Result(); orv != ov || ore != oe {
		t.Errorf(`Try(badFunc).Result() should equal (Try(badFunc).Value(), Try(badFunc).Err()); got (%v, %v != %v, %v)`, orv, ore, ov, oe)
	}
	if oes, exp := out.Error(), ot; oes != exp {
		t.Errorf(`Try(badFunc).Error() = %q, want %q`, oes, exp)
	}
	info = out.info
	if len(info) != 0 {
		t.Errorf(`len(Try(badFunc).Info()) = %q, want %q`, len(info), 0)
	}
}
