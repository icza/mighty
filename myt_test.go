package mighty

import (
	"errors"
	"math"
	"testing"
)

type TBMock struct {
	testing.TB

	errCalls int
}

func (m *TBMock) Errorf(format string, args ...interface{}) {
	m.errCalls++
}

func TestMytDeqEqNeq(t *testing.T) {
	tb := &TBMock{}
	myt := Myt{tb}

	errs := []error{errors.New("test error")}
	cases := []struct {
		exp, got                             interface{}
		errs                                 []error
		eqErrCalls, deqErrCalls, neqErrCalls int
	}{
		{1, 1, nil, 0, 0, 1},
		{1, 2, nil, 1, 1, 0},
		{1, "3", nil, 2, 2, 0},
		{1, 1, errs, 1, 1, 1},
		{1, 2, errs, 1, 1, 1},
		{1, "3", errs, 2, 2, 1},
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
		tb.errCalls = 0
		myt.Deq(c.exp, c.got, c.errs...)
		if c.deqErrCalls != tb.errCalls {
			t.Errorf("[i=%d] Expected: %d, got: %d", i, c.deqErrCalls, tb.errCalls)
		}
	}
}

func TestMytDeq(t *testing.T) {
	tb := &TBMock{}
	myt := Myt{tb}

	errs := []error{errors.New("test error")}
	cases := []struct {
		exp, got interface{}
		errs     []error
		errCalls int
	}{
		{[]int{1, 2}, []int{1, 2}, nil, 0},
		{[]int{1, 2}, []int{2, 2}, nil, 1},
		{[]int{1, 2}, "x", nil, 2},
		{[]int{1, 2}, []int{1, 2}, errs, 1},
		{[]int{1, 2}, []int{2, 2}, errs, 1},
		{[]int{1, 2}, "x", errs, 2},
	}

	for i, c := range cases {
		tb.errCalls = 0
		myt.Deq(c.exp, c.got, c.errs...)
		if c.errCalls != tb.errCalls {
			t.Errorf("[i=%d] Expected: %d, got: %d", i, c.errCalls, tb.errCalls)
		}
	}
}

func TestMytNear(t *testing.T) {
	tb := &TBMock{}
	myt := Myt{tb}

	errs := []error{errors.New("test error")}
	cases := []struct {
		exp, got, eps float64
		errs          []error
		errCalls      int
	}{
		{1.0, 1.0, 1e-6, nil, 0},
		{1.0, 1.0, 1e-6, errs, 1},
		{1.0, 1.001, 1e-2, nil, 0},
		{1.0, 1.001, 1e-2, errs, 1},
		{1.0, 1.001, 1e-4, nil, 1},
		{1.0, 1.001, 1e-4, errs, 1},
	}

	for i, c := range cases {
		exp, got := c.exp, c.got
		for j := 0; j < 2; j++ { // 2 cycles to switch got end exp
			if j > 0 {
				exp, got = got, exp
			}
			tb.errCalls = 0
			myt.Near(c.exp, c.got, c.eps, c.errs...)
			if c.errCalls != tb.errCalls {
				t.Errorf("[i=%d] Expected: %d, got: %d", i, c.errCalls, tb.errCalls)
			}
		}
	}
}

func TestFuncs(t *testing.T) {
	tb := &TBMock{}

	Eq(tb)(1, 2)
	Deq(tb)(1, 2)
	Neq(tb)(1, 2)
	Near(tb)(1, 1, 1e-6)

	eq, neq := EqNeq(tb)
	eq(1, 2)
	neq(1, 2)

	ExpEq(tb)(1)(2, nil)
	ExpDeq(tb)(1)(2, nil)
	ExpNeq(tb)(1)(2, nil)
	ExpNear(tb)(1, 1e-6)(1, nil)

	eq, expEq := EqExpEq(tb)
	eq(1, 2)
	expEq(1)(2, nil)
}

func TestNearLogic(t *testing.T) {
	inf, neginf, nan := math.Inf(1), math.Inf(-1), math.NaN()
	cases := []struct {
		a, b, eps float64
		exp       bool
	}{
		{1.0, 1.0, 1e-6, true},
		{1.0, 1.001, 1e-6, false},
		{1.0, 1.001, 1e-2, true},

		// Corner cases
		{inf, 1.001, 1e-2, false},
		{neginf, 1.001, 1e-2, false},
		{inf, inf, 1e-2, true},
		{neginf, neginf, 1e-2, true},
		{inf, neginf, 1e-2, false},

		{1.0, 1.1, inf, true},
		{1.0, inf, inf, false},
		{inf, inf, inf, true},
		{neginf, neginf, inf, true},

		{1.0, nan, 1e10, false},
		{1.0, nan, inf, false},
		{nan, nan, 1e10, false},
		{nan, nan, inf, false},

		{1.0, 1.0, nan, true},
		{1.0, 1.001, nan, false},
		{inf, inf, nan, true},
		{neginf, neginf, nan, true},
	}

	for i, c := range cases {
		if got := NearLogic(c.a, c.b, c.eps); c.exp != got {
			t.Errorf("[i=%d] Expected: %v, got: %v", i, c.exp, got)
		}
	}
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
