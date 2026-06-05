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
	cases := []struct {
		a, b int
		want int
	}{
		{48, 18, 6},
		{10, 5, 5},
		{7, 3, 1},
		{0, 5, 5},
	}
	for _, c := range cases {
		if got := GCD(c.a, c.b); got != c.want {
			t.Errorf("GCD(%d, %d) = %d; want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLCM(t *testing.T) {
	cases := []struct {
		a, b int
		want int
	}{
		{4, 6, 12},
		{5, 7, 35},
		{0, 5, 0},
	}
	for _, c := range cases {
		if got := LCM(c.a, c.b); got != c.want {
			t.Errorf("LCM(%d, %d) = %d; want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestFibonacci(t *testing.T) {
	cases := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{10, 55},
	}
	for _, c := range cases {
		if got := Fibonacci(c.n); got != c.want {
			t.Errorf("Fibonacci(%d) = %d; want %d", c.n, got, c.want)
		}
	}
}

func TestReverse(t *testing.T) {
	if got := Reverse("gopher"); got != "rehpog" {
		t.Errorf("Reverse(\"gopher\") = %s; want \"rehpog\"", got)
	}
	if got := Reverse("🚀🚀"); got != "🚀🚀" {
		t.Errorf("Reverse(\"🚀🚀\") = %s; want \"🚀🚀\"", got)
	}
}

func TestSubstring(t *testing.T) {
	cases := []struct {
		s          string
		start, end int
		want       string
	}{
		{"hello", 1, 3, "hel"},
		{"hello", 2, 4, "ell"},
		{"hello", -1, 2, "he"},
		{"hello", 4, 10, "lo"},
		{"hello", 4, 2, ""},
		{"🚀rocket", 1, 1, "🚀"},
		{"🚀rocket", 2, 7, "rocket"},
	}
	for _, c := range cases {
		if got := Substring(c.s, c.start, c.end); got != c.want {
			t.Errorf("Substring(%q, %d, %d) = %q; want %q", c.s, c.start, c.end, got, c.want)
		}
	}
}