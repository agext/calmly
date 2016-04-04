// Copyright 2016 ALRUX Inc.
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
/*
Package calmly implements convenient runtime panic recovery and handling


When a panic condition needs to be handled by the program (rather than crashing it), wrap the code that can trigger such condition in a `Try`, which allows you to `Catch` the panic for further processing.

The `Outcome` of a `Try`ed code also offers convenience methods to:
- `KeepCalm` downgrading a panic to an error condition;
- `Escalate` upgrading a panic to a fatal error;
- `Log` the error, panic or fatal condition, using the appropriate logger method - presumably triggering a new panic or exiting the program.
*/
package calmly

import (
	"fmt"
	"runtime"
)

// Outcome represents the state of a `Try`ed call, including information about
// any panic it may have triggered, as well as the returned value and error, if applicable.
type Outcome struct {
	val   interface{}
	err   error
	level int8
	code  int
	text  string
	info  []string
}

// Try calls the function it receives as argument, recovering from any panic it may cause
func Try(f interface{}) (o *Outcome) {
	defer func() {
		if err := recover(); err != nil {
			o.level, o.code, o.text = PANIC, ERR_TRY_PANIC, fmt.Sprintf("panic: %s", err)
			o.addInfo(2, "debug.stack")
		}
	}()

	o = &Outcome{level: OK}
	switch f := f.(type) {
	case func():
		f()
	case func() error:
		o.err = f()
	case func() interface{}:
		o.val = f()
	case func() (interface{}, error):
		o.val, o.err = f()
	default:
		o = &Outcome{
			level: ERROR,
			code:  ERR_TRY_ARG,
			text:  fmt.Sprintf("Try: unsupported argument type %T", f),
		}
	}
	return
}

// Catch calls the provided function passing the receiver Outcome as argument,
// only if the Outcome is at PANIC level.
func (this *Outcome) Catch(f func(*Outcome)) *Outcome {
	if this.level == PANIC {
		f(this)
	}
	return this
}

// KeepCalm downgrades a PANIC to ERROR level, to avoid triggering a panic upon
// logging the outcome.
func (this *Outcome) KeepCalm() *Outcome {
	if this.level == PANIC {
		this.level = ERROR
	}
	return this
}

// Escalate converts a PANIC into a FATAL condition, to trigger program
// termination upon logging the outcome.
func (this *Outcome) Escalate() *Outcome {
	if this.level == PANIC {
		this.level = FATAL
	}
	return this
}

// Log sends the error-condition Outcome to the provided log, using the appropriate
// logging function: FATAL conditions are logged using Fatal(), PANIC using
// Panic(), and ERROR using Print(). Non-error conditions are not logged
// because there is no information stored in the Outcome, beside
// what the Try-ed function returned (and is better suited to log itself).
func (this *Outcome) Log(log Logger) *Outcome {
	switch this.level {
	case FATAL:
		log.Fatal(this)
	case PANIC:
		log.Panic(this)
	case ERROR:
		log.Print(this)
	}
	return this
}

// Level returns the error level stored by the receiver.
func (this *Outcome) Level() int8 {
	return this.level
}

// SetLevel sets the error level stored by the receiver.
func (this *Outcome) SetLevel(l int8) *Outcome {
	if levelName(l) != "?" {
		this.level = l
	}
	return this
}

// Code returns the error code stored by the receiver.
func (this *Outcome) Code() int {
	return this.code
}

// SetCode sets the error code stored by the receiver.
func (this *Outcome) SetCode(c int) *Outcome {
	this.code = c
	return this
}

// Text returns the error text stored by the receiver.
func (this *Outcome) Text() string {
	return this.text
}

// SetText sets the error text stored by the receiver.
func (this *Outcome) SetText(t string) *Outcome {
	this.text = t
	return this
}

// Info returns the error info stored by the receiver.
func (this *Outcome) Info() []string {
	return this.info
}

// addInfo adds (more) error info to the receiver.
func (this *Outcome) addInfo(calldepth int, s ...string) *Outcome {
	for i, line := range s {
		if line == "debug.stack" {
			calldepth *= 2
			buffer := make([]byte, 4096)
			buffer = buffer[:runtime.Stack(buffer, true)]
			var p1, p2, l int
			for j, c := range buffer {
				if c == 10 {
					if l == 0 {
						p1 = j + 1
					} else if l == calldepth {
						p2 = j + 1
						break
					}
					l++
				}
			}
			if p2 > 0 {
				s[i] = string(buffer[:p1]) + string(buffer[p2:])
			} else {
				s[i] = string(buffer)
			}
			break
		}
	}
	this.info = append(this.info, s...)
	return this
}

// AddInfo adds (more) error info to the receiver.
func (this *Outcome) AddInfo(s ...string) *Outcome {
	return this.addInfo(2, s...)
}

// Value provides the value returned by the Try-ed function, if any.
func (this *Outcome) Value() interface{} {
	return this.val
}

// Err provides the error returned by the Try-ed function, if any.
func (this *Outcome) Err() error {
	return this.err
}

// Result provides the value and error returned by the Try-ed function, if any.
func (this *Outcome) Result() (interface{}, error) {
	return this.val, this.err
}

// Error returns a string representation of the Outcome if it is in an error condition,
// or an empty string if no error or panic occurred. Note that the Try-ed function
// returning a non-nil error does not constitute an error condition for the Outcome.
// That error value can be retrieved via Err or Result.
// This is also useful for satisfying the `error` interface.
func (this *Outcome) Error() string {
	if this.level == OK {
		return ""
	}
	if this.code != 0 {
		return this.text + fmt.Sprintf(" (code: 0x%04x)", this.code)
	}
	return this.text
}
