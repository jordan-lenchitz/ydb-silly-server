package main

import (
	"fmt"
	"strings"
)

// --- LOG Utilities ---

func LogWrite(lvl, msg string) {
	if conn == nil {
		fmt.Printf("[%s] (NO DB) %s\n", lvl, msg)
		return
	}
	node := conn.Node("^YDBLOGS", lvl)
	id := node.Incr(1)
	
	conn.Node("^YDBLOGS", lvl, id, "ts").Set(GetDateTime())
	conn.Node("^YDBLOGS", lvl, id, "msg").Set(msg)
	conn.Node("^YDBLOGS", lvl, id, "job").Set(fmt.Sprintf("%d", startTime.Unix())) // Using startTime as a pseudo-job ID or similar
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
	if conn == nil {
		return ""
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
	if conn == nil {
		return ""
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
	if conn == nil {
		return
	}
	if lvl == "" {
		conn.Node("^YDBLOGS").Kill()
	} else {
		conn.Node("^YDBLOGS", lvl).Kill()
	}
}

// --- DATA Utilities ---

func GetGlobal(gbl string, subs ...any) string {
	if conn == nil {
		return ""
	}
	return conn.Node("^"+gbl, subs...).Get()
}

func SetGlobal(gbl string, val string, subs ...any) {
	if conn == nil {
		return
	}
	conn.Node("^"+gbl, subs...).Set(val)
}

func KillGlobal(gbl string, subs ...any) {
	if conn == nil {
		return
	}
	conn.Node("^"+gbl, subs...).Kill()
}

func GlobalExists(gbl string, subs ...any) bool {
	if conn == nil {
		return false
	}
	node := conn.Node("^"+gbl, subs...)
	return !node.HasNone()
}
