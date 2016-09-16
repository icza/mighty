package mighty

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var packageName = reflect.TypeOf(Myt{}).PkgPath() // e.g. "github.com/icza/mighty"

// Myt is a wrapper around a testing.TB value which is usually
// a *testing.T or a *testing.B, arming it with short utility methods.
type Myt struct {
	testing.TB
}

// Eq reports an error if exp != got, or an optional non-nil error is provided.
func (m Myt) Eq(exp, got interface{}, errs ...error) {
	m.ExpEq(exp)(got, errs...)
}

// Neq reports an error if exp == got, or an optional non-nil error is provided.
func (m Myt) Neq(v1, v2 interface{}, errs ...error) {
	m.ExpNeq(v1)(v2, errs...)
}

// ExpEq takes the expected value and returns a function which
// only takes the 'got' value and an optional error.
//
// The following multiline code:
//     exp := some_expected_value
//     got, err := SomeFunc()
//     Eq(exp, got, err)
//
// Is equivalent to this single line:
//     ExpEq(exp)(SoomeFunc())
func (m Myt) ExpEq(exp interface{}) func(got interface{}, errs ...error) {
	return func(got interface{}, errs ...error) {
		var err error
		if len(errs) > 0 && errs[0] != nil {
			err = errs[0]
		}
		if exp != got || err != nil {
			function, line := getFuncLine()
			if err == nil {
				m.Errorf("[%s:%d] Expected: %v, got: %v", function, line, exp, got)
			} else {
				m.Errorf("[%s:%d] Expected: %v, got: %v, error: %v", function, line, exp, got, err)
			}
			// Common mistake is to provide constants as exp whose default value will be applied
			// when packed into interface{} which might not be the case in case of direct comparison.
			// Provide warning for such likely cause.
			if exp != got && exp != nil && got != nil {
				if texp, tgot := reflect.TypeOf(exp), reflect.TypeOf(got); texp != tgot {
					m.Errorf("\tType of expected and got does not match! exp type: %v, got type: %v", texp, tgot)
				}
			}
		}
	}
}

// ExpNeq takes one value and returns a function which
// takes only the 2nd value and an optional error.
//
// The following multiline code:
//     v1 := some_value1
//     v2, err := SomeFunc()
//     Neq(v1, v2, err)
//
// Is equivalent to this single line:
//     ExpNeq(v1)(SoomeFunc())
func (m Myt) ExpNeq(v1 interface{}) func(v2 interface{}, errs ...error) {
	return func(v2 interface{}, errs ...error) {
		var err error
		if len(errs) > 0 && errs[0] != nil {
			err = errs[0]
		}
		if v1 == v2 || err != nil {
			function, line := getFuncLine()
			if err == nil {
				m.Errorf("[%s:%d] Expected mismatch: %v, got: %v", function, line, v1, v2)
			} else {
				m.Errorf("[%s:%d] Expected mismatch: %v, got: %v, error: %v", function, line, v1, v2, err)
			}
		}
	}
}

// getFuncLine reports the function name and line number of the first caller
// that is not from this package.
func getFuncLine() (function string, line int) {
	callers := make([]uintptr, 20)
	count := runtime.Callers(1, callers)
	frames := runtime.CallersFrames(callers[:count])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if !strings.HasPrefix(frame.Function, packageName) {
			return frame.Function, frame.Line
		}
	}

	return "<unknown_file>", -1
}

// Func2ArgsErr is a type describing a function which takes 2 interface{} arguments
// and an optional error (in the form of a variadic parameter).
type Func2ArgsErr func(interface{}, interface{}, ...error)

// Func1ArgFunc1ArgErr is a type describing a function which takes 1 interface{} argument
// and returns a function which takes 1 interface{} argument and an optional error (in the form of a variadic parameter).
type Func1ArgFunc1ArgErr func(interface{}) func(interface{}, ...error)

// Eq returns a method value of Myt{t}.Eq.
// tb may be a *testing.T or *testing.B value.
func Eq(tb testing.TB) Func2ArgsErr {
	return Myt{tb}.Eq
}

// Neq returns a method value of Myt{t}.Neq.
// tb may be a *testing.T or *testing.B value.
func Neq(tb testing.TB) Func2ArgsErr {
	return Myt{tb}.Neq
}

// EqNeq returns 2 method values: Myt{t}.Eq and Myt{t}.Neq.
// tb may be a *testing.T or *testing.B value.
func EqNeq(tb testing.TB) (Func2ArgsErr, Func2ArgsErr) {
	myt := Myt{tb}
	return myt.Eq, myt.Neq
}

// ExpEq returns a method value of Myt{t}.ExpEq.
// tb may be a *testing.T or *testing.B value.
func ExpEq(tb testing.TB) Func1ArgFunc1ArgErr {
	return Myt{tb}.ExpEq
}

// ExpNeq returns a method value of Myt{t}.ExpNeq.
// tb may be a *testing.T or *testing.B value.
func ExpNeq(tb testing.TB) Func1ArgFunc1ArgErr {
	return Myt{tb}.ExpNeq
}

// EqExpEq returns 2 method values: Myt{t}.Eq and Myt{t}.ExpEq.
// tb may be a *testing.T or *testing.B value.
func EqExpEq(tb testing.TB) (Func2ArgsErr, Func1ArgFunc1ArgErr) {
	myt := Myt{tb}
	return myt.Eq, myt.ExpEq
}
