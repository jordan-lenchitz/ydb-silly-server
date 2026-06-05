package main

import (
	"testing"
)

func TestFactorial(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{0, 1},
		{1, 1},
		{5, 120},
	}
	for _, c := range cases {
		got := Factorial(c.n)
		if got != c.want {
			t.Errorf("Factorial(%d) == %d, want %d", c.n, got, c.want)
		}
	}
}

func TestIsPrime(t *testing.T) {
	if !IsPrime(2) {
		t.Error("2 should be prime")
	}
	if IsPrime(4) {
		t.Error("4 should not be prime")
	}
	if !IsPrime(17) {
		t.Error("17 should be prime")
	}
}

func TestGCD(t *testing.T) {
	if got := GCD(48, 18); got != 6 {
		t.Errorf("GCD(48, 18) = %d; want 6", got)
	}
}

func TestReverse(t *testing.T) {
	if got := Reverse("gopher"); got != "rehpog" {
		t.Errorf("Reverse(\"gopher\") = %s; want \"rehpog\"", got)
	}
}