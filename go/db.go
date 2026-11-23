package main

import (
	"fmt"
	"strings"
	"sync"
)

var dbMu sync.Mutex

func IsDBConnected() bool {
	dbMu.Lock()
	defer dbMu.Unlock()
	return conn != nil
}

func normalizeGlobalName(gbl string) string {
	gbl = strings.TrimSpace(gbl)
	gbl = strings.TrimPrefix(gbl, "^")
	return gbl
}

// --- LOG Utilities ---

func LogWrite(lvl, msg string) {
	lvl = strings.TrimSpace(lvl)
	if lvl == "" {
		lvl = "INFO"
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		fmt.Printf("[%s] (NO DB) %s\n", lvl, msg)
		return
	}
	if m != nil {
		if _, err := m.CallErr("LogWriteM", lvl, msg); err == nil {
			return
		}
	}
	node := conn.Node("^YDBLOGS", lvl)
	id := node.Incr(1)

	conn.Node("^YDBLOGS", lvl, id, "ts").Set(GetDateTime())
	conn.Node("^YDBLOGS", lvl, id, "msg").Set(msg)
	conn.Node("^YDBLOGS", lvl, id, "job").Set(fmt.Sprintf("%d", startTime.Unix()))
}

func LogInfo(msg string) {
	LogWrite("INFO", msg)
}

func LogWarn(msg string) {
	LogWrite("WARN", msg)
}

func LogError(msg string) {
	LogWrite("ERROR", msg)
}

func ReadLogs(lvl string, count int) string {
	lvl = strings.TrimSpace(lvl)
	if lvl == "" || count <= 0 {
		return ""
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return ""
	}
	if m != nil {
		if v, err := m.CallErr("LogReadM", lvl, count); err == nil && v != nil {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	idVal := conn.Node("^YDBLOGS", lvl).Get()
	if idVal == "" {
		return ""
	}

	var id int
	fmt.Sscanf(idVal, "%d", &id)

	var res strings.Builder
	start := id - count + 1
	if start < 1 {
		start = 1
	}

	for i := id; i >= start; i-- {
		ts := conn.Node("^YDBLOGS", lvl, i, "ts").Get()
		msg := conn.Node("^YDBLOGS", lvl, i, "msg").Get()
		res.WriteString(fmt.Sprintf("%s [%s] %s\n", lvl, ts, msg))
	}

	return res.String()
}

func GetLogStats() string {
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return ""
	}
	if m != nil {
		if v, err := m.CallErr("LogStatsM"); err == nil && v != nil {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	var res strings.Builder
	node := conn.Node("^YDBLOGS")

	for _, lvl := range node.Children() {
		count := conn.Node("^YDBLOGS", lvl).Get()
		res.WriteString(fmt.Sprintf("%s:%s;", lvl, count))
	}

	return res.String()
}

func ClearLogs(lvl string) {
	lvl = strings.TrimSpace(lvl)
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return
	}
	if m != nil {
		if _, err := m.CallErr("LogClearM", lvl); err == nil {
			return
		}
	}
	if lvl == "" {
		conn.Node("^YDBLOGS").Kill()
	} else {
		conn.Node("^YDBLOGS", lvl).Kill()
	}
}

// --- DATA Utilities ---

func GetGlobal(gbl string, subs ...any) string {
	gbl = normalizeGlobalName(gbl)
	if gbl == "" {
		return ""
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return ""
	}
	return conn.Node("^"+gbl, subs...).Get()
}

func SetGlobal(gbl string, val string, subs ...any) {
	gbl = normalizeGlobalName(gbl)
	if gbl == "" {
		return
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return
	}
	conn.Node("^"+gbl, subs...).Set(val)
}

func KillGlobal(gbl string, subs ...any) {
	gbl = normalizeGlobalName(gbl)
	if gbl == "" {
		return
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return
	}
	conn.Node("^"+gbl, subs...).Kill()
}

func GlobalExists(gbl string, subs ...any) bool {
	gbl = normalizeGlobalName(gbl)
	if gbl == "" {
		return false
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return false
	}
	node := conn.Node("^"+gbl, subs...)
	return !node.HasNone()
}
