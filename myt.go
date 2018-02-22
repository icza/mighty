package mighty

import (
	"fmt"
	"math"
	"path/filepath"
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

// Deq reports an error if !reflect.DeepEqual(exp, got), or an optional
// non-nil error is provided.
func (m Myt) Deq(exp, got interface{}, errs ...error) {
	m.ExpDeq(exp)(got, errs...)
}

// Neq reports an error if exp == got, or an optional non-nil error is provided.
func (m Myt) Neq(v1, v2 interface{}, errs ...error) {
	m.ExpNeq(v1)(v2, errs...)
}

// Near reports an error if the float64 exp is not "near" to got,
// or an optional non-nil error is provided.
// "near" is defined by the NearLogic function variable.
func (m Myt) Near(exp, got, eps float64, errs ...error) {
	m.ExpNear(exp, eps)(got, errs...)
}

// ExpEq takes the expected value and returns a function which
// only takes the 'got' value and an optional error.
//
// The following multiline code:
//     got, err := SomeFunc()
//     Eq(exp, got, err)
//
// Is equivalent to this single line:
//     ExpEq(exp)(SomeFunc())
func (m Myt) ExpEq(exp interface{}) func(got interface{}, errs ...error) {
	return m.expEqDeq(exp, false)
}

// ExpDeq takes the expected value and returns a function which
// only takes the 'got' value and an optional error.
//
// The following multiline code:
//     got, err := SomeFunc()
//     Deq(exp, got, err)
//
// Is equivalent to this single line:
//     ExpDeq(exp)(SomeFunc())
func (m Myt) ExpDeq(exp interface{}) func(got interface{}, errs ...error) {
	return m.expEqDeq(exp, true)
}

// expEqDeq takes the expected value and returns a function which
// only takes the 'got' value and an optional error.
// Whether deep equality has to be used is controled by the deep argument.
func (m Myt) expEqDeq(exp interface{}, deep bool) func(got interface{}, errs ...error) {
	return func(got interface{}, errs ...error) {
		err := getErr(errs...)
		var eq bool
		var separator string
		if deep {
			eq = reflect.DeepEqual(exp, got)
			// When checking deep equality, values may be "big" (complex),
			// so align 'got' under 'Expected' for easy visual comparibility.
			separator = "\n\t    "
		} else {
			eq = exp == got
		}
		if eq && err == nil {
			return
		}

		if err == nil {
			m.Errorf("%s\n\tExpected: %v,%s got: %v", getFuncLine(), exp, separator, got)
		} else {
			m.Errorf("%s\n\tExpected: %v,%s got: %v, error: %v", getFuncLine(), exp, separator, got, err)
		}
		// Common mistake is to provide constants as exp whose default value will be applied
		// when packed into interface{} which might not be the case in case of direct comparison.
		// Provide warning for such likely cause.
		if !eq && exp != nil && got != nil {
			if texp, tgot := reflect.TypeOf(exp), reflect.TypeOf(got); texp != tgot {
				m.Errorf("\tTypes of expected and got do not match! exp type: %v, got type: %v", texp, tgot)
			}
		}
	}
}

// ExpNeq takes one value and returns a function which
// takes only the 2nd value and an optional error.
//
// The following multiline code:
//     v2, err := SomeFunc()
//     Neq(v1, v2, err)
//
// Is equivalent to this single line:
//     ExpNeq(v1)(SomeFunc())
func (m Myt) ExpNeq(v1 interface{}) func(v2 interface{}, errs ...error) {
	return func(v2 interface{}, errs ...error) {
		err := getErr(errs...)
		if v1 != v2 && err == nil {
			return
		}

		if err == nil {
			m.Errorf("%s\n\tExpected mismatch: %v, got: %v", getFuncLine(), v1, v2)
		} else {
			m.Errorf("%s\n\tExpected mismatch: %v, got: %v, error: %v", getFuncLine(), v1, v2, err)
		}
	}
}

// ExpNear takes the expected and epsilon values and returns a function which
// only takes the 'got' value and an optional error.
//
// The following multiline code:
//     got, err := SomeFunc()
//     Near(exp, got, eps, err)
//
// Is equivalent to this single line:
//     ExpNear(exp, eps)(SomeFunc())
func (m Myt) ExpNear(exp, eps float64) func(got float64, errs ...error) {
	return func(got float64, errs ...error) {
		err := getErr(errs...)
		if err == nil && NearLogic(exp, got, eps) {
			return
		}

		if err == nil {
			m.Errorf("%s\n\tExpected: %v, got: %v, with eps: %v)",
				getFuncLine(), exp, got, eps)
		} else {
			m.Errorf("%s\n\tExpected: %v, got: %v, with eps: %v, error: %v", getFuncLine(), exp, got, eps, err)
		}
	}
}

// NearLogic is a variable holding a function which is responsible to
// decide if 2 float64 numbers are near to each other (given an epsilon).
// It is used by the Myt.Near() and Myt.ExpNear() functions.
// Default value is NearFunc, but you may set your own function.
var NearLogic = NearFunc

// NearFunc checks if 2 float64 numbers are "near" to each other.
// The caller is responsible to provide a sensible epsilon.
// This is the default NearLogic, but you may set your own function.
//
// "near" is defined as the following:
//     near := Math.Abs(a - b) < eps
//
// Corner cases:
//  1. if a==b, result is true (eps will not be checked, may be NaN)
//  2. Inf is near to Inf (even if eps=NaN; consequence of 1.)
//  3. -Inf is near to -Inf (even if eps=NaN; consequence of 1.)
//  4. NaN is not near to anything (not even to NaN)
//  5. eps=Inf results in true (unless any of a or b is NaN)
func NearFunc(a, b, eps float64) bool {
	// Quick check, also handles infinities:
	if a == b {
		return true
	}

	return math.Abs(a-b) < eps
}

// getErr is a utility function which returns the optional error
// from the variadic errors.
func getErr(errs ...error) error {
	if len(errs) > 0 && errs[0] != nil {
		return errs[0]
	}
	return nil
}

// getFuncLine returns a formatted string containing the function name,
// file name and line number of the first caller that is not from this package.
func getFuncLine() string {
	var function, file string
	var line int

	callers := make([]uintptr, 20)
	count := runtime.Callers(1, callers)
	frames := runtime.CallersFrames(callers[:count])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if !strings.HasPrefix(frame.Function, packageName) {
			function, file, line = frame.Function, filepath.Base(frame.File), frame.Line
			break
		}
	}

	if function == "" {
		function, file, line = "<unknown_func>", "<unknown_file>", -1
	}

	return fmt.Sprintf("Func: %s, File: %s:%d", function, file, line)
}

// Func2ArgsErr is a type describing a function which takes 2 interface{} arguments
// and an optional error (in the form of a variadic parameter).
type Func2ArgsErr func(interface{}, interface{}, ...error)

// Func1ArgFunc1ArgErr is a type describing a function which takes 1 interface{} argument
// and returns a function which takes 1 interface{} argument and an optional error (in the form of a variadic parameter).
type Func1ArgFunc1ArgErr func(interface{}) func(interface{}, ...error)

// Func3FloatsErr is a type describing a function which takes 3 float64 arguments
// and an optional error (in the form of a variadic parameter).
type Func3FloatsErr func(float64, float64, float64, ...error)

// Func2FloatsFunc1FloatErr is a type describing a function which takes 2 float64 arguments
// and returns a function which takes 1 float64 argument and an optional error (in the form of a variadic parameter).
type Func2FloatsFunc1FloatErr func(float64, float64) func(float64, ...error)

// Eq returns a method value of Myt{tb}.Eq.
// tb may be a *testing.T or *testing.B value.
func Eq(tb testing.TB) Func2ArgsErr {
	return Myt{tb}.Eq
}

// Deq returns a method value of Myt{tb}.Deq.
// tb may be a *testing.T or *testing.B value.
func Deq(tb testing.TB) Func2ArgsErr {
	return Myt{tb}.Deq
}

// Neq returns a method value of Myt{tb}.Neq.
// tb may be a *testing.T or *testing.B value.
func Neq(tb testing.TB) Func2ArgsErr {
	return Myt{tb}.Neq
}

// Near returns a method value of Myt{tb}.Near.
// tb may be a *testing.T or *testing.B value.
func Near(tb testing.TB) Func3FloatsErr {
	return Myt{tb}.Near
}

// EqNeq returns 2 method values: Myt{tb}.Eq and Myt{tb}.Neq.
// tb may be a *testing.T or *testing.B value.
func EqNeq(tb testing.TB) (Func2ArgsErr, Func2ArgsErr) {
	myt := Myt{tb}
	return myt.Eq, myt.Neq
}

// ExpEq returns a method value of Myt{tb}.ExpEq.
// tb may be a *testing.T or *testing.B value.
func ExpEq(tb testing.TB) Func1ArgFunc1ArgErr {
	return Myt{tb}.ExpEq
}

// ExpDeq returns a method value of Myt{tb}.ExpDeq.
// tb may be a *testing.T or *testing.B value.
func ExpDeq(tb testing.TB) Func1ArgFunc1ArgErr {
	return Myt{tb}.ExpDeq
}

// ExpNeq returns a method value of Myt{tb}.ExpNeq.
// tb may be a *testing.T or *testing.B value.
func ExpNeq(tb testing.TB) Func1ArgFunc1ArgErr {
	return Myt{tb}.ExpNeq
}

// ExpNear returns a method value of Myt{tb}.ExpNear.
// tb may be a *testing.T or *testing.B value.
func ExpNear(tb testing.TB) Func2FloatsFunc1FloatErr {
	return Myt{tb}.ExpNear
}

// EqExpEq returns 2 method values: Myt{tb}.Eq and Myt{tb}.ExpEq.
// tb may be a *testing.T or *testing.B value.
func EqExpEq(tb testing.TB) (Func2ArgsErr, Func1ArgFunc1ArgErr) {
	myt := Myt{tb}
	return myt.Eq, myt.ExpEq
}
