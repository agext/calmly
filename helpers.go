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

// Outcome levels match logging levels in agext/log
const (
	OK int8 = 0

	ERROR int8 = iota + 3
	PANIC
	FATAL
)

// Predefined error codes
const (
	ERR_TRY_ARG int = iota
	ERR_TRY_PANIC
)

func levelName(l int8) string {
	switch l {
	case OK:
		return "OK"
	case ERROR:
		return "ERROR"
	case PANIC:
		return "PANIC"
	case FATAL:
		return "FATAL"
	}
	return "?"
}

// Logger defines the interface expected by the Log method of Outcome
type Logger interface {
	Fatal(...interface{})
	Panic(...interface{})
	Print(...interface{})
}
