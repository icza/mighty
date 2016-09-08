package mighty

import (
	"runtime"
	"testing"
)

// Myt is a wrapper around *testing.T, arming it with short utility methods.
type Myt struct {
	*testing.T
}

// Eq checks if exp and got are equal, and if not, reports it as an error.
func (m Myt) Eq(exp, got interface{}) {
	if exp != got {
		_, file, line, _ := runtime.Caller(2)
		m.T.Errorf("[%s:%d] Expected: %v, got: %v", file, line, exp, got)
	}
}

// Neq checks if v1 and v2 are not equal, but if they are, reports it as an error.
func (m Myt) Neq(v1, v2 interface{}) {
	if v1 == v2 {
		_, file, line, _ := runtime.Caller(2)
		m.T.Errorf("[%s:%d] Expected mismatch: %v, got: %v", file, line, v1, v2)
	}
}

// Eq returns a method value of Myt{t}.Eq.
func Eq(t *testing.T) func(interface{}, interface{}) {
	return Myt{t}.Eq
}

// Neq returns a method value of Myt{t}.Neq.
func Neq(t *testing.T) func(interface{}, interface{}) {
	return Myt{t}.Neq
}

// EqNeq returns 2 method values: Myt{t}.Eq and Myt{t}.Neq.
func EqNeq(t *testing.T) (func(interface{}, interface{}), func(interface{}, interface{})) {
	myt := Myt{t}
	return myt.Eq, myt.Neq
}
