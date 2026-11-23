package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// --- WEB Utilities ---

func URLEncode(s string) string {
	return url.QueryEscape(s)
}

func URLDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

func HexToDec(h string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(h), 16, 64)
}

func ParseHeaders(h string) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(h, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if key != "" {
				headers[key] = val
			}
		}
	}
	return headers
}

// --- MUMPS Call-In & XEXE Utilities ---

func CallM(fn string, args ...any) (any, error) {
	dbMu.Lock()
	defer dbMu.Unlock()
	if m == nil {
		return nil, fmt.Errorf("engine not connected")
	}
	return m.CallErr(fn, args...)
}

func XExe(code string) (int, error) {
	if !IsDBConnected() {
		return 0, fmt.Errorf("engine not connected")
	}
	jobID := os.Getpid()
	KillGlobal("XOUT", jobID)

	res, err := CallM("XExe", code)
	if err != nil {
		return 0, err
	}

	return res.(int), nil
}

func EvalM(expr string) (string, error) {
	res, err := CallM("Eval", expr)
	if err != nil {
		return "", err
	}
	if s, ok := res.(string); ok {
		return s, nil
	}
	return fmt.Sprintf("%v", res), nil
}
