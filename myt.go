package mighty

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var packageName = reflect.TypeOf(Myt{}).PkgPath() // e.g. "github.com/icza/mighty"

// Myt is a wrapper around *testing.T, arming it with short utility methods.
type Myt struct {
	*testing.T
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
			file, line := getFileLine()
			if err == nil {
				m.T.Errorf("[%s:%d] Expected: %v, got: %v", file, line, exp, got)
			} else {
				m.T.Errorf("[%s:%d] Expected: %v, got: %v, error: %v", file, line, exp, got, err)
			}
			// Common mistake is to provide constants as exp whose default value will be applied
			// when packed into interface{} which might not be the case in case of direct comparision.
			// Provide warning for such likely cause.
			if exp != got && exp != nil && got != nil {
				if texp, tgot := reflect.TypeOf(exp), reflect.TypeOf(got); texp != tgot {
					m.T.Errorf("\tType of expected and got does not match! exp type: %v, got type: %v", texp, tgot)
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
			file, line := getFileLine()
			if err == nil {
				m.T.Errorf("[%s:%d] Expected mismatch: %v, got: %v", file, line, v1, v2)
			} else {
				m.T.Errorf("[%s:%d] Expected mismatch: %v, got: %v, error: %v", file, line, v1, v2, err)
			}
		}
	}
}

// getFileLine reports the file name and line number of the first caller
// that is not from this package.
func getFileLine() (file string, line int) {
	callers := make([]uintptr, 20)
	count := runtime.Callers(1, callers)
	for i := 0; i < count; i++ {
		pc := callers[i]
		if fd := runtime.FuncForPC(pc); !strings.HasPrefix(fd.Name(), packageName) {
			file, line = fd.FileLine(pc)
			line-- // TODO: line is actual line +1, WHY??
			return
		}
	}
	return "<unknown_file>", -1
}

// Eq returns a method value of Myt{t}.Eq.
func Eq(t *testing.T) func(interface{}, interface{}, ...error) {
	return Myt{t}.Eq
}

// Neq returns a method value of Myt{t}.Neq.
func Neq(t *testing.T) func(interface{}, interface{}, ...error) {
	return Myt{t}.Neq
}

// EqNeq returns 2 method values: Myt{t}.Eq and Myt{t}.Neq.
func EqNeq(t *testing.T) (func(interface{}, interface{}, ...error), func(interface{}, interface{}, ...error)) {
	myt := Myt{t}
	return myt.Eq, myt.Neq
}

// ExpEq returns a method value of Myt{t}.ExpEq.
func ExpEq(t *testing.T) func(interface{}) func(interface{}, ...error) {
	return Myt{t}.ExpEq
}

// ExpNeq returns a method value of Myt{t}.ExpNeq.
func ExpNeq(t *testing.T) func(interface{}) func(interface{}, ...error) {
	return Myt{t}.ExpNeq
}
