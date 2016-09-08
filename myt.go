package mighty

import (
	"runtime"
	"testing"
)

// Myt is a wrapper around *testing.T, arming it with short utility methods.
type Myt struct {
	*testing.T
}

// Eq reports an error if exp != got, or an optional non-nil error is provided.
func (m Myt) Eq(exp, got interface{}, errs ...error) {
	var err error
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
	}
	if exp != got || err != nil {
		_, file, line, _ := runtime.Caller(2)
		if err == nil {
			m.T.Errorf("[%s:%d] Expected: %v, got: %v", file, line, exp, got)
		} else {
			m.T.Errorf("[%s:%d] Expected: %v, got: %v, error: %v", file, line, exp, got, err)
		}
	}
}

// Neq reports an error if exp == got, or an optional non-nil error is provided.
func (m Myt) Neq(v1, v2 interface{}, errs ...error) {
	var err error
	if len(errs) > 0 && errs[0] != nil {
		err = errs[0]
	}
	if v1 == v2 {
		_, file, line, _ := runtime.Caller(2)
		if err == nil {
			m.T.Errorf("[%s:%d] Expected mismatch: %v, got: %v", file, line, v1, v2)
		} else {
			m.T.Errorf("[%s:%d] Expected mismatch: %v, got: %v, error: %v", file, line, v1, v2, err)
		}
	}
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
