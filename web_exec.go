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
	return strconv.ParseInt(h, 16, 64)
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
			headers[key] = val
		}
	}
	return headers
}

// --- XEXE Utilities ---

func XExe(code string) (int, error) {
	jobID := os.Getpid()
	KillGlobal("XOUT", jobID)
	
	// Since the codebase is being refactored to Go, and M code execution is still required,
	// we still need a way to call XECUTE in YottaDB.
	// We will use the existing 'm' (MFunctions) for this until we fully remove it,
	// OR we can implement a Go function that calls a generic M shim.
	
	// For now, let's assume we still have access to the 'm' instance in main.go
	// or we use conn.MCall if we move the import.
	
	res, err := m.CallErr("XExe", code)
	if err != nil {
		return 0, err
	}
	
	return res.(int), nil
}
