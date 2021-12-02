// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/rglujing/gf.

package gerror

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rglujing/gf/errors/gcode"
	"github.com/rglujing/gf/internal/utils"
	"io"
	"runtime"
	"strings"
)

// Error is custom error for additional features.
type Error struct {
	error error      // Wrapped error.
	stack stack      // Stack array, which records the stack information when this error is created or wrapped.
	text  string     // Error text, which is created by New* functions.
	code  gcode.Code // Error code if necessary.
}

const (
	// Filtering key for current error module paths.
	stackFilterKeyLocal = "/errors/gerror/gerror"
)

var (
	// goRootForFilter is used for stack filtering purpose.
	// Mainly for development environment.
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

// Error implements the interface of Error, it returns all the error as string.
func (err *Error) Error() string {
	if err == nil {
		return ""
	}
	errStr := err.text
	if errStr == "" && err.code != nil {
		errStr = err.code.Message()
	}
	if err.error != nil {
		if errStr != "" {
			errStr += ": "
		}
		errStr += err.error.Error()
	}
	return errStr
}

// Code returns the error code.
// It returns CodeNil if it has no error code.
func (err *Error) Code() gcode.Code {
	if err == nil {
		return gcode.CodeNil
	}
	return err.code
}

// Cause returns the root cause error.
func (err *Error) Cause() error {
	if err == nil {
		return nil
	}
	loop := err
	for loop != nil {
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				// Internal Error struct.
				loop = e
			} else if e, ok := loop.error.(apiCause); ok {
				// Other Error that implements ApiCause interface.
				return e.Cause()
			} else {
				return loop.error
			}
		} else {
			// return loop
			// To be compatible with Case of https://github.com/pkg/errors.
			return errors.New(loop.text)
		}
	}
	return nil
}

// Format formats the frame according to the fmt.Formatter interface.
//
// %v, %s   : Print all the error string;
// %-v, %-s : Print current level error string;
// %+s      : Print full stack error list;
// %+v      : Print the error string and full stack error list;
func (err *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('-'):
			if err.text != "" {
				io.WriteString(s, err.text)
			} else {
				io.WriteString(s, err.Error())
			}
		case s.Flag('+'):
			if verb == 's' {
				io.WriteString(s, err.Stack())
			} else {
				io.WriteString(s, err.Error()+"\n"+err.Stack())
			}
		default:
			io.WriteString(s, err.Error())
		}
	}
}

// Stack returns the stack callers as string.
// It returns an empty string if the <err> does not support stacks.
func (err *Error) Stack() string {
	if err == nil {
		return ""
	}
	var (
		loop   = err
		index  = 1
		buffer = bytes.NewBuffer(nil)
	)
	for loop != nil {
		buffer.WriteString(fmt.Sprintf("%d. %-v\n", index, loop))
		index++
		formatSubStack(loop.stack, buffer)
		if loop.error != nil {
			if e, ok := loop.error.(*Error); ok {
				loop = e
			} else {
				buffer.WriteString(fmt.Sprintf("%d. %s\n", index, loop.error.Error()))
				index++
				break
			}
		} else {
			break
		}
	}
	return buffer.String()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func (err *Error) Current() error {
	if err == nil {
		return nil
	}
	return &Error{
		error: nil,
		stack: err.stack,
		text:  err.text,
		code:  err.code,
	}
}

// Next returns the next level error.
// It returns nil if current level error or the next level error is nil.
func (err *Error) Next() error {
	if err == nil {
		return nil
	}
	return err.error
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// Note that do not use pointer as its receiver here.
func (err *Error) MarshalJSON() ([]byte, error) {
	return []byte(`"` + err.Error() + `"`), nil
}

// formatSubStack formats the stack for error.
func formatSubStack(st stack, buffer *bytes.Buffer) {
	if st == nil {
		return
	}
	index := 1
	space := "  "
	for _, p := range st {
		if fn := runtime.FuncForPC(p - 1); fn != nil {
			file, line := fn.FileLine(p - 1)
			// Custom filtering.
			if !utils.IsDebugEnabled() {
				if strings.Contains(file, utils.StackFilterKeyForGoFrame) {
					continue
				}
			} else {
				if strings.Contains(file, stackFilterKeyLocal) {
					continue
				}
			}
			// Avoid stack string like "<autogenerated>"
			if strings.Contains(file, "<") {
				continue
			}
			// Ignore GO ROOT paths.
			if goRootForFilter != "" &&
				len(file) >= len(goRootForFilter) &&
				file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			// Graceful indent.
			if index > 9 {
				space = " "
			}
			buffer.WriteString(fmt.Sprintf(
				"   %d).%s%s\n    \t%s:%d\n",
				index, space, fn.Name(), file, line,
			))
			index++
		}
	}
}
