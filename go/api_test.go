package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"lang.yottadb.com/go/yottadb/v2"
)

func TestMain(mRun *testing.M) {
	gin.SetMode(gin.TestMode)
	if os.Getenv("ydb_dist") != "" || os.Getenv("YDB_DIST") != "" {
		if err := tryInitYDB(); err != nil {
			println("[TestMain] tryInitYDB warning:", err.Error())
		}
	}
	code := mRun.Run()
	if db != nil {
		yottadb.Shutdown(db)
	}
	os.Exit(code)
}

func TestVMStatusEndpoint(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/vm/status", nil)
	req.Header.Set("X-VM-ID", "test-vm-nonexistent-123")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", w.Code)
	}

	var response VMState
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status == "" {
		t.Error("Expected Status in response, got empty string")
	}

	if response.System == nil {
		t.Fatal("Expected System map in response, got nil")
	}
	if _, ok := response.System["engine"]; !ok {
		t.Error("Expected 'engine' field in system metadata")
	}
	if _, ok := response.System["uptime"]; !ok {
		t.Error("Expected 'uptime' field in system metadata")
	}
}

func TestVMProvisionAndLifecycle(t *testing.T) {
	r := setupRouter()

	body := bytes.NewBufferString(`{"id":"ydb-test-vm-42","cores":7,"sync":true}`)
	req, _ := http.NewRequest("POST", "/api/vm/provision", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("POST /api/vm/provision returned %d: %s", w.Code, w.Body.String())
	}

	var provResp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &provResp); err != nil {
		t.Fatalf("Failed to parse provision response: %v", err)
	}
	if provResp["instanceId"] != "ydb-test-vm-42" {
		t.Errorf("Expected instanceId ydb-test-vm-42, got %v", provResp["instanceId"])
	}

	if IsDBConnected() {
		statusReq, _ := http.NewRequest("GET", "/api/vm/status", nil)
		statusReq.Header.Set("X-VM-ID", "ydb-test-vm-42")
		sw := httptest.NewRecorder()
		r.ServeHTTP(sw, statusReq)

		var st VMState
		if err := json.Unmarshal(sw.Body.Bytes(), &st); err != nil {
			t.Fatalf("Failed to unmarshal status: %v", err)
		}
		if st.Status != "RUNNING" {
			t.Errorf("Expected VM status RUNNING after sync provision, got %s", st.Status)
		}
		if !st.Provisioned {
			t.Error("Expected Provisioned=true")
		}
		if st.Metadata["cpuLimit"] != "7.0" {
			t.Errorf("Expected cpuLimit 7.0, got %v", st.Metadata["cpuLimit"])
		}

		// List VMs
		listReq, _ := http.NewRequest("GET", "/api/vm/list", nil)
		lw := httptest.NewRecorder()
		r.ServeHTTP(lw, listReq)
		if lw.Code != http.StatusOK || !strings.Contains(lw.Body.String(), "ydb-test-vm-42") {
			t.Errorf("Expected /api/vm/list to contain ydb-test-vm-42, got %s", lw.Body.String())
		}

		// Delete VM
		delReq, _ := http.NewRequest("DELETE", "/api/vm/ydb-test-vm-42", nil)
		dw := httptest.NewRecorder()
		r.ServeHTTP(dw, delReq)
		if dw.Code != http.StatusOK {
			t.Errorf("Expected DELETE /api/vm/ydb-test-vm-42 200, got %d", dw.Code)
		}
	}
}

func TestExecuteNilSafety(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("POST", "/api/execute", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for nil body, got %v", w.Code)
	}

	emptyReq, _ := http.NewRequest("POST", "/execute", bytes.NewBufferString(`{"mCode":"   "}`))
	emptyReq.Header.Set("Content-Type", "application/json")
	ew := httptest.NewRecorder()
	r.ServeHTTP(ew, emptyReq)

	if ew.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for empty mCode, got %v", ew.Code)
	}
}

func TestExecuteMUMPSCode(t *testing.T) {
	if !IsDBConnected() {
		t.Skip("YottaDB not connected; skipping live XEXE HTTP execution test")
	}
	r := setupRouter()

	cases := []struct {
		name       string
		mCode      string
		wantSubstr string
	}{
		{"SingleLineNoNewline", `W "BYTE-PERFECT-100"`, "BYTE-PERFECT-100"},
		{"MultiLineOutput", `W "LINE1",!,"LINE2",!`, "LINE1\nLINE2"},
		{"YottaDBIntrinsicZV", `W $ZYRELEASE," ",$ZV`, "YottaDB"},
		{"RuntimeErrorTrapped", `S X=1/0`, "[MUMPS-ERROR]"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			payload, _ := json.Marshal(map[string]string{"mCode": tc.mCode})
			req, _ := http.NewRequest("POST", "/api/execute", bytes.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Expected 200 OK, got %d: %s", w.Code, w.Body.String())
			}
			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Invalid JSON: %v", err)
			}
			out, _ := resp["output"].(string)
			if !strings.Contains(out, tc.wantSubstr) {
				t.Errorf("Output %q did not contain %q", out, tc.wantSubstr)
			}
		})
	}
}

func TestGlobalCRUDEndpoints(t *testing.T) {
	r := setupRouter()

	// Bad JSON test
	badReq, _ := http.NewRequest("POST", "/api/global/TESTGBL", bytes.NewBufferString(`{invalid`))
	badReq.Header.Set("Content-Type", "application/json")
	bw := httptest.NewRecorder()
	r.ServeHTTP(bw, badReq)
	if bw.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for bad JSON on POST /api/global/:name, got %d", bw.Code)
	}

	// Set Global
	setBody := bytes.NewBufferString(`{"subscripts":["users","42","role"],"value":"gopher-mumps"}`)
	setReq, _ := http.NewRequest("POST", "/api/global/^TESTGBL", setBody)
	setReq.Header.Set("Content-Type", "application/json")
	sw := httptest.NewRecorder()
	r.ServeHTTP(sw, setReq)
	if sw.Code != http.StatusOK {
		t.Fatalf("Expected 200 on POST /api/global/^TESTGBL, got %d", sw.Code)
	}

	// Get Global
	getReq, _ := http.NewRequest("GET", "/api/global/TESTGBL?subs=users,42,role", nil)
	gw := httptest.NewRecorder()
	r.ServeHTTP(gw, getReq)
	if gw.Code != http.StatusOK {
		t.Fatalf("Expected 200 on GET /api/global/TESTGBL, got %d", gw.Code)
	}

	if IsDBConnected() {
		var getResp map[string]interface{}
		_ = json.Unmarshal(gw.Body.Bytes(), &getResp)
		if getResp["value"] != "gopher-mumps" {
			t.Errorf("Expected global value 'gopher-mumps', got %v", getResp["value"])
		}

		// Inspect tree
		treeReq, _ := http.NewRequest("GET", "/api/global/TESTGBL/tree?max=10", nil)
		tw := httptest.NewRecorder()
		r.ServeHTTP(tw, treeReq)
		if tw.Code != http.StatusOK || !strings.Contains(tw.Body.String(), "gopher-mumps") {
			t.Errorf("Expected tree dump to contain gopher-mumps, got %s", tw.Body.String())
		}

		// Delete Global
		delReq, _ := http.NewRequest("DELETE", "/api/global/TESTGBL", nil)
		dw := httptest.NewRecorder()
		r.ServeHTTP(dw, delReq)
		if dw.Code != http.StatusOK {
			t.Errorf("Expected 200 on DELETE /api/global/TESTGBL, got %d", dw.Code)
		}
	}
}

func TestFilesEndpointSuffixSorting(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/files", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from /api/files, got %d", w.Code)
	}

	var resp struct {
		Files []struct {
			Name   string `json:"name"`
			Folder string `json:"folder"`
			Path   string `json:"path"`
			Suffix string `json:"suffix"`
		} `json:"files"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse /api/files JSON: %v", err)
	}
	if len(resp.Files) == 0 {
		t.Fatal("Expected non-empty files list from /api/files")
	}

	for i := 1; i < len(resp.Files); i++ {
		prev, curr := resp.Files[i-1], resp.Files[i]
		if prev.Suffix > curr.Suffix {
			t.Errorf("Files not sorted by suffix at index %d: %q > %q", i, prev.Suffix, curr.Suffix)
		}
		if prev.Suffix == curr.Suffix && prev.Name > curr.Name {
			t.Errorf("Files with identical suffix %q not sorted by name: %q > %q", prev.Suffix, prev.Name, curr.Name)
		}
	}
}

func TestDocsAndCORSEndpoints(t *testing.T) {
	r := setupRouter()

	for _, route := range []string{"/api/docs", "/v1beta/api", "/api/health", "/api/sys/intrinsics", "/api/wisdom", "/api/dogs", "/api/frogs", "/api/cats"} {
		req, _ := http.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("GET %s returned %d, want 200", route, w.Code)
		}
		if w.Body.Len() == 0 {
			t.Errorf("GET %s returned empty body", route)
		}
	}

	// Preflight OPTIONS test
	optReq, _ := http.NewRequest("OPTIONS", "/api/execute", nil)
	ow := httptest.NewRecorder()
	r.ServeHTTP(ow, optReq)
	if ow.Code != http.StatusOK {
		t.Errorf("OPTIONS /api/execute returned %d, want 200", ow.Code)
	}
	if ow.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Missing CORS header Access-Control-Allow-Origin")
	}
}

func TestLogsForgeCatsAndMUMPSTestEndpoints(t *testing.T) {
	r := setupRouter()

	LogInfo("api test log entry")
	logReq, _ := http.NewRequest("GET", "/api/logs?level=INFO&count=5", nil)
	lw := httptest.NewRecorder()
	r.ServeHTTP(lw, logReq)
	if lw.Code != http.StatusOK {
		t.Errorf("GET /api/logs returned %d", lw.Code)
	}

	if IsDBConnected() {
		// POST /api/cats
		catReq, _ := http.NewRequest("POST", "/api/cats", bytes.NewBufferString(`{"name":"MUMPS-CAT","mood":"INTRINSIC"}`))
		catReq.Header.Set("Content-Type", "application/json")
		cw := httptest.NewRecorder()
		r.ServeHTTP(cw, catReq)
		if cw.Code != http.StatusOK {
			t.Errorf("POST /api/cats returned %d: %s", cw.Code, cw.Body.String())
		}

		// POST /api/forge & GET /api/decode/A
		forgeReq, _ := http.NewRequest("POST", "/api/forge", nil)
		fw := httptest.NewRecorder()
		r.ServeHTTP(fw, forgeReq)
		if fw.Code != http.StatusOK {
			t.Errorf("POST /api/forge returned %d: %s", fw.Code, fw.Body.String())
		}

		decReq, _ := http.NewRequest("GET", "/api/decode/A", nil)
		dw := httptest.NewRecorder()
		r.ServeHTTP(dw, decReq)
		if dw.Code != http.StatusOK || !strings.Contains(dw.Body.String(), "ASCII_HEX_41") {
			t.Errorf("GET /api/decode/A unexpected body: %s", dw.Body.String())
		}

		// GET /api/mumps/test
		mtReq, _ := http.NewRequest("GET", "/api/mumps/test", nil)
		mw := httptest.NewRecorder()
		r.ServeHTTP(mw, mtReq)
		if mw.Code != http.StatusOK || !strings.Contains(mw.Body.String(), `"status":"PASS"`) {
			t.Errorf("GET /api/mumps/test failed: %s", mw.Body.String())
		}
	}
}
