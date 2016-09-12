package mighty

import (
	"errors"
	"testing"
)

type TBMock struct {
	testing.TB

	errCalls int
}

func (m *TBMock) Errorf(format string, args ...interface{}) {
	m.errCalls++
}

func TestMytEqNeq(t *testing.T) {
	tb := &TBMock{}
	myt := Myt{tb}

	errs := []error{errors.New("test error")}
	cases := []struct {
		exp, got                interface{}
		errs                    []error
		eqErrCalls, neqErrCalls int
	}{
		{1, 1, nil, 0, 1},
		{1, 2, nil, 1, 0},
		{1, "3", nil, 2, 0},
		{1, 1, errs, 1, 1},
		{1, 2, errs, 1, 1},
		{1, "3", errs, 2, 1},
	}

	for i, c := range cases {
		tb.errCalls = 0
		myt.Eq(c.exp, c.got, c.errs...)
		if c.eqErrCalls != tb.errCalls {
			t.Errorf("[i=%d] Expected: %d, got: %d", i, c.eqErrCalls, tb.errCalls)
		}
		tb.errCalls = 0
		myt.Neq(c.exp, c.got, c.errs...)
		if c.neqErrCalls != tb.errCalls {
			t.Errorf("[i=%d] Expected: %d, got: %d", i, c.neqErrCalls, tb.errCalls)
		}
	}
}

func TestFuncs(t *testing.T) {
	tb := &TBMock{}

	Eq(tb)(1, 2)
	Neq(tb)(1, 2)

	eq, neq := EqNeq(tb)
	eq(1, 2)
	neq(1, 2)

	ExpEq(tb)(1)(2, nil)
	ExpNeq(tb)(1)(2, nil)

	eq, expEq := EqExpEq(tb)
	eq(1, 2)
	expEq(1)(2, nil)
}

func TestGetFileLineUnknown(t *testing.T) {
	// We need a "deep" stack
	var f func(int)

	f = func(n int) {
		if n < 25 {
			f(n + 1)
		} else {
			Myt{&TBMock{}}.Eq(1, 2)
		}
	}
	f(0)
}
