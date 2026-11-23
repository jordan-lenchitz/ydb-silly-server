package main

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestConcurrentDatabaseAccess(t *testing.T) {
	const workers = 50
	const iterations = 100
	var wg sync.WaitGroup

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				if j%25 == 0 {
					LogInfo("stress test log")
				}
				SetGlobal("STRESS", "val", id, j)
				_ = GetGlobal("STRESS", id, j)
				_ = GlobalExists("STRESS", id, j)
			}
		}(i)
	}
	wg.Wait()
	KillGlobal("STRESS")
}

func TestLogInvariant(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogWrite panicked: %v", r)
		}
	}()

	LogWrite("", "")
	LogWrite("CRITICAL", "Something happened")
	LogInfo("info message")
	LogWarn("warn message")
	LogError("error message")

	if IsDBConnected() {
		logs := ReadLogs("CRITICAL", 5)
		if !strings.Contains(logs, "Something happened") {
			t.Errorf("Expected CRITICAL log to contain 'Something happened', got %q", logs)
		}
		stats := GetLogStats()
		if !strings.Contains(stats, "CRITICAL:") {
			t.Errorf("Expected GetLogStats() to include CRITICAL:, got %q", stats)
		}
		ClearLogs("CRITICAL")
		if after := ReadLogs("CRITICAL", 5); after != "" {
			t.Errorf("Expected empty logs after ClearLogs, got %q", after)
		}
	}
}

func TestGlobalPersistenceLifecycle(t *testing.T) {
	KillGlobal("UNITGBL")
	if GlobalExists("UNITGBL", "k1") {
		t.Error("Expected UNITGBL(k1) not to exist after KillGlobal")
	}

	SetGlobal("^UNITGBL", "v1", "k1", "sub1")
	if IsDBConnected() {
		if !GlobalExists("UNITGBL", "k1", "sub1") {
			t.Error("Expected UNITGBL(k1,sub1) to exist after SetGlobal")
		}
		if got := GetGlobal("^UNITGBL", "k1", "sub1"); got != "v1" {
			t.Errorf("GetGlobal(^UNITGBL, k1, sub1) = %q, want v1", got)
		}
		KillGlobal("UNITGBL", "k1", "sub1")
		if GlobalExists("UNITGBL", "k1", "sub1") {
			t.Error("Expected UNITGBL(k1,sub1) to be removed after KillGlobal")
		}
	}

	// Empty global name safety
	SetGlobal("", "ignored")
	if got := GetGlobal(""); got != "" {
		t.Errorf("Expected empty string for empty global name, got %q", got)
	}
	if GlobalExists("") {
		t.Error("Expected GlobalExists(\"\") == false")
	}
	KillGlobal("")
}

func TestVMMgrOperations(t *testing.T) {
	vmID := fmt.Sprintf("vm-unit-%d", 9001)
	CleanVM(vmID)

	res := Provision(vmID, "5")
	if res != "SUCCESS" {
		t.Fatalf("Provision returned %q, want SUCCESS", res)
	}

	if IsDBConnected() {
		if st := GetVM(vmID, "status"); st != "PROVISIONING" {
			t.Errorf("Expected PROVISIONING status, got %q", st)
		}
		if cpu := GetMeta(vmID, "cpuLimit"); cpu != "5.0" {
			t.Errorf("Expected cpuLimit 5.0, got %q", cpu)
		}
		Complete(vmID)
		if st := GetVM(vmID, "status"); st != "RUNNING" {
			t.Errorf("Expected RUNNING status after Complete, got %q", st)
		}
		StopVM(vmID)
		if st := GetVM(vmID, "status"); st != "STOPPED" {
			t.Errorf("Expected STOPPED status after StopVM, got %q", st)
		}
		CleanVM(vmID)
		if st := GetVM(vmID, "status"); st != "" {
			t.Errorf("Expected empty status after CleanVM, got %q", st)
		}
	}
}
