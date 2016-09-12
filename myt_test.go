package mighty

import (
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

	cases := []struct {
		exp, got                interface{}
		eqErrCalls, neqErrCalls int
	}{
		{1, 1, 0, 1},
		{1, 2, 1, 0},
		{1, "3", 2, 0},
	}

	for i, c := range cases {
		tb.errCalls = 0
		myt.Eq(c.exp, c.got)
		if c.eqErrCalls != tb.errCalls {
			t.Errorf("[i=%d] Expected: %d, got: %d", i, c.eqErrCalls, tb.errCalls)
		}
		tb.errCalls = 0
		myt.Neq(c.exp, c.got)
		if c.neqErrCalls != tb.errCalls {
			t.Errorf("[i=%d] Expected: %d, got: %d", i, c.neqErrCalls, tb.errCalls)
		}
	}
}
