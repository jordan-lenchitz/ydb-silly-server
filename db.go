package main

import (
	"fmt"
	"strings"

	"lang.yottadb.com/go/yottadb/v2"
)

// --- LOG Utilities ---

func LogWrite(lvl, msg string) {
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

func ClearLogs(lvl string) {
	if lvl == "" {
		conn.Node("^YDBLOGS").Kill()
	} else {
		conn.Node("^YDBLOGS", lvl).Kill()
	}
}

func GetLogStats() string {
	var res strings.Builder
	node := conn.Node("^YDBLOGS")
	
	curr, err := node.SubscriptNext("")
	for err == nil && curr != nil {
		lvl := curr.(string)
		count := conn.Node("^YDBLOGS", lvl).Get()
		res.WriteString(fmt.Sprintf("%s:%s;", lvl, count))
		curr, err = node.SubscriptNext(lvl)
	}
	
	return res.String()
}

// --- DATA Utilities ---

func GetGlobal(gbl string, subs ...any) string {
	return conn.Node("^"+gbl, subs...).Get()
}

func SetGlobal(gbl string, val string, subs ...any) {
	conn.Node("^"+gbl, subs...).Set(val)
}

func KillGlobal(gbl string, subs ...any) {
	conn.Node("^"+gbl, subs...).Kill()
}

func GlobalExists(gbl string, subs ...any) bool {
	// yottadb-go doesn't have a direct $D equivalent that returns 0-11,
	// but Get() returning empty might mean not exists, or we can use SubscriptNext.
	// Actually, let's just check if Get returns something or has children.
	val := conn.Node("^"+gbl, subs...).Get()
	if val != "" {
		return true
	}
	// Check for children
	next, _ := conn.Node("^"+gbl, subs...).SubscriptNext("")
	return next != nil
}
