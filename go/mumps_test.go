package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNativeMUMPSSuite(t *testing.T) {
	if !IsDBConnected() {
		t.Skip("YottaDB not connected; skipping native MUMPS suite")
	}

	res, err := CallM("RunMTests")
	if err != nil {
		t.Fatalf("CallM(RunMTests) failed: %v", err)
	}
	summary, _ := res.(string)
	t.Logf("MUMPS ZTEST Summary: %s", summary)

	if !strings.Contains(summary, "FAIL:0|") {
		totalStr := GetGlobal("YDBTEST", "TOTAL")
		t.Fatalf("Native MUMPS test suite reported failures (%s, total=%s)", summary, totalStr)
	}
}

func TestMUMPSEvalAndIntrinsics(t *testing.T) {
	if !IsDBConnected() {
		t.Skip("YottaDB not connected; skipping MUMPS Eval & Intrinsics test")
	}

	cases := []struct {
		name string
		expr string
		want string
	}{
		{"MathFact7", "$$FACT^MATH(7)", "5040"},
		{"MathIsPrime97", "$$ISPRIME^MATH(97)", "1"},
		{"MathGCD", "$$GCD^MATH(1071,462)", "21"},
		{"MathLCM", "$$LCM^MATH(21,6)", "42"},
		{"MathFib12", "$$FIB^MATH(12)", "144"},
		{"MathBitAnd", "$$BITAND^MATH(240,170)", "160"},
		{"MathBitOr", "$$BITOR^MATH(240,15)", "255"},
		{"MathBitXor", "$$BITXOR^MATH(255,170)", "85"},
		{"MathBitCnt", "$$BITCNT^MATH(255)", "8"},
		{"StrRevUTF8", `$$REV^STR("MUMPS🚀")`, "🚀SPMUM"},
		{"StrUp", `$$UP^STR("yottadb")`, "YOTTADB"},
		{"StrLow", `$$LOW^STR("YOTTADB")`, "yottadb"},
		{"StrBLenEmoji", `$$BLEN^STR("💯")`, "4"},
		{"StrLenEmoji", `$$LEN^STR("💯")`, "1"},
		{"WebURLRoundtrip", `$$URLDEC^WEB($$URLENC^WEB("yottadb 🚀 💯"))`, "yottadb 🚀 💯"},
		{"SysHealth", `$$HEALTH^SYS()`, "OK"},
		{"EvalErrorTrap", `1/0`, "[EVAL-ERROR]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := EvalM(tc.expr)
			if err != nil {
				t.Fatalf("EvalM(%s) failed: %v", tc.expr, err)
			}
			if !strings.Contains(got, tc.want) {
				t.Errorf("EvalM(%s) = %q, want substring %q", tc.expr, got, tc.want)
			}
		})
	}
}

func TestOfflineFallbackResilience(t *testing.T) {
	// Temporarily simulate nil conn/m to ensure 100% panic-free fallback paths
	dbMu.Lock()
	savedConn, savedM := conn, m
	conn, m = nil, nil
	dbMu.Unlock()

	defer func() {
		dbMu.Lock()
		conn, m = savedConn, savedM
		dbMu.Unlock()
	}()

	if IsDBConnected() {
		t.Fatal("Expected IsDBConnected() == false when conn is nil")
	}

	// DB & VM operations should safely no-op or return empty values
	SetGlobal("OFFLINE", "val", "sub")
	if got := GetGlobal("OFFLINE", "sub"); got != "" {
		t.Errorf("Expected empty GetGlobal while offline, got %q", got)
	}
	if GlobalExists("OFFLINE", "sub") {
		t.Error("Expected GlobalExists == false while offline")
	}
	KillGlobal("OFFLINE", "sub")
	LogInfo("offline log test")
	if logs := ReadLogs("INFO", 5); logs != "" {
		t.Errorf("Expected empty ReadLogs while offline, got %q", logs)
	}
	if stats := GetLogStats(); stats != "" {
		t.Errorf("Expected empty GetLogStats while offline, got %q", stats)
	}
	ClearLogs("INFO")
	if vms := ListVMs(); len(vms) != 0 {
		t.Errorf("Expected empty ListVMs while offline, got %v", vms)
	}

	r := setupRouter()
	for _, route := range []string{"/api/vm/status", "/api/cats", "/api/wisdom", "/api/health"} {
		req, _ := http.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Offline GET %s returned %d, want 200", route, w.Code)
		}
	}
}

func BenchmarkMUMPSEval(b *testing.B) {
	if !IsDBConnected() {
		b.Skip("YottaDB not connected")
	}
	for i := 0; i < b.N; i++ {
		_, _ = EvalM("$$FIB^MATH(10)")
	}
}
