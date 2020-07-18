package user

import (
	"testing"
)

type TestCase struct {
	ID    uint32
	IsOdd bool
}

func TestIsOdd(t *testing.T) {
	cases := []*TestCase{
		{1, false},
		{2, true},
	}

	for _, tc := range cases {
		u := NewUser(tc.ID, "test")
		if u.IsOdd() != tc.IsOdd {
			t.Errorf("unexpected result for %v, have %v, want %v", u.ID, u.IsOdd(), tc.IsOdd)
		}
	}
}
