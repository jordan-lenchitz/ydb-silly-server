package main

import (
	"sync"
	"testing"
)

func TestConcurrentDatabaseAccess(t *testing.T) {
	
	const workers = 100
	const iterations = 1000
	var wg sync.WaitGroup
	
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				LogInfo("stress test log")
				_ = GetGlobal("STRESS", id, j)
				SetGlobal("STRESS", "val", id, j)
				_ = GlobalExists("STRESS", id, j)
				
				if j%100 == 0 {
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestLogInvariant(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("LogWrite panicked: %v", r)
		}
	}()
	
	LogWrite("", "")
	LogWrite("CRITICAL", "Something happened")
}