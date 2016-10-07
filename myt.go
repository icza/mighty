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

// Near reports an error if the float64 exp is not "near" to got,
// or an optional non-nil error is provided.
// "near" is defined by the NearLogic() function.
func (m Myt) Near(exp, got, eps float64, errs ...error) {
	m.ExpNear(exp, eps)(got, errs...)
}

// ExpEq takes the expected value and returns a function which
// only takes the 'got' value and an optional error.
//
// The following multiline code:
//     got, err := SomeFunc()
//     Eq(someExpectedValue, got, err)
//
// Is equivalent to this single line:
//     ExpEq(someExpectedValue)(SomeFunc())
func (m Myt) ExpEq(exp interface{}) func(got interface{}, errs ...error) {
	return func(got interface{}, errs ...error) {
		err := getErr(errs...)
		if exp == got && err == nil {
			return
		}

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

// ExpNeq takes one value and returns a function which
// takes only the 2nd value and an optional error.
//
// The following multiline code:
//     v2, err := SomeFunc()
//     Neq(someValue1, v2, err)
//
// Is equivalent to this single line:
//     ExpNeq(someValue1)(SomeFunc())
func (m Myt) ExpNeq(v1 interface{}) func(v2 interface{}, errs ...error) {
	return func(v2 interface{}, errs ...error) {
		err := getErr(errs...)
		if v1 != v2 && err == nil {
			return
		}

		function, line := getFuncLine()
		if err == nil {
			m.Errorf("[%s:%d] Expected mismatch: %v, got: %v", function, line, v1, v2)
		} else {
			m.Errorf("[%s:%d] Expected mismatch: %v, got: %v, error: %v", function, line, v1, v2, err)
		}
	}
}

// ExpNear takes the expected and epslion values and returns a function which
// only takes the 'got' value and an optional error.
//
// The following multiline code:
//     got, err := SomeFunc()
//     Near(someExpectedValue, got, someEpsilon, err)
//
// Is equivalent to this single line:
//     ExpNear(someExpectedValue, someEpsilon)(SomeFunc())
func (m Myt) ExpNear(exp, eps float64) func(got float64, errs ...error) {
	return func(got float64, errs ...error) {
		err := getErr(errs...)
		if err == nil && NearLogic(exp, got, eps) {
			return
		}

		function, line := getFuncLine()
		if err == nil {
			m.Errorf("[%s:%d] Expected: %v, got: %v, with eps: %v)", function, line, exp, got, eps)
		} else {
			m.Errorf("[%s:%d] Expected: %v, got: %v, with eps: %v, error: %v", function, line, exp, got, eps, err)
		}
	}
}

// NearLogic checks if 2 float64 numbers are "near" to each other.
// The caller is responsible to provide a sensible epsilon.
// "near" is defined as the following:
//     near := Math.Abs(a - b) < eps
//
// Corner cases:
//  1. if a==b, result is true (eps will not be checked, may be NaN)
//  2. Inf is near to Inf (even if eps=NaN; consequence of 1.)
//  3. -Inf is near to -Inf (even if eps=NaN; consequence of 1.)
//  4. NaN is not near to anything (not even to NaN)
//  5. eps=Inf results in true (unless any of a or b is NaN)
func NearLogic(a, b, eps float64) bool {
	// Quick check, also handles infinities:
	if a == b {
		return true
	}

	diff := a - b
	if diff < 0 {
		diff = -diff
	}

	return diff < eps
}

// getErr is a utility function which returns the optional error
// from the variadic errors.
func getErr(errs ...error) error {
	if len(errs) > 0 && errs[0] != nil {
		return errs[0]
	}
	return nil
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
