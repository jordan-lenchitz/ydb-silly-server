package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// --- SYS Utilities ---

var startTime = time.Now()

func GetUptime() string {
	return fmt.Sprintf("%.0f", time.Since(startTime).Seconds())
}

func GetDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// --- MATH Utilities ---

func Factorial(n int) int {
	res := 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	return res
}

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 {
		return true
	}
	if n%2 == 0 {
		return false
	}
	for i := 3; i <= int(math.Sqrt(float64(n))); i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return (a * b) / GCD(a, b)
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func Pow(b, e int) int {
	return int(math.Pow(float64(b), float64(e)))
}

func Round(n float64, d int) float64 {
	p := math.Pow(10, float64(d))
	return math.Round(n*p) / p
}

func Fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// --- STR Utilities ---

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Upper(s string) string {
	return strings.ToUpper(s)
}

func Lower(s string) string {
	return strings.ToLower(s)
}

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func Substring(s string, start, end int) string {
	if start < 1 {
		start = 1
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		return ""
	}
	return s[start-1 : end]
}
