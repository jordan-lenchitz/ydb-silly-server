package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"lang.yottadb.com/go/yottadb/v2"
)

/**
 * 🚀 YOTTADB UNIRONIC BACKEND V4.0.0 (GO NATIVE)
 * Pure Gopher logic, YottaDB Persistence & Native MUMPS Call-Ins.
 */

const (
	BaseDir = "."
)

var (
	folders = []string{"mumps", "go", "js", "md", "logs"}
	db      *yottadb.DB
	conn    *yottadb.Conn
	m       *yottadb.MFunctions
)

// Call-In table bridging Go to native MUMPS routines in mumps/*.m
const mCallTable = `
	XExe: int XEXE^XEXE(string)
	Eval: string[4096] EVAL^XEXE(string)
	Version: string[512] VERSION^XEXE()
	Intrinsics: string[4096] INTRINSICS^XEXE()
	SysMetrics: string[4096] METRICS^SYS()
	SysHealth: string[256] HEALTH^SYS()
	SysDateTime: string[256] DATETIME^SYS()
	VMProvision: string[256] PROVISION^VMMGR(string, string)
	VMComplete: string[256] COMPLETE^VMMGR(string)
	VMStop: string[256] STOP^VMMGR(string)
	VMGet: string[1024] GETVM^VMMGR(string, string)
	VMMeta: string[1024] GETMETA^VMMGR(string, string)
	VMList: string[4096] LIST^VMMGR()
	VMCount: int COUNT^VMMGR()
	VMClean: string[256] CLEAN^VMMGR(string)
	LogWriteM: int WRITE^LOG(string, string)
	LogReadM: string[8192] READ^LOG(string, int)
	LogStatsM: string[2048] STATS^LOG()
	LogClearM: string[256] CLEAR^LOG(string)
	DataCount: int COUNT^DATA(string)
	DataDump: string[16384] DUMP^DATA(string, int)
	CatsSeed: int SEED^ZDB()
	CatsAdd: int ADD^ZDB(string, string)
	CatsJSON: string[8192] C^ZDB()
	ForgeEmoji: int FORGE^ZFORGE()
	DecodeEmoji: string[2048] DECODE^ZFORGE(string)
	RigidRoute: string[8192] ROUTE^ZREQ(string)
	RunMTests: string[4096] ALL^ZTEST()
`

type VMState struct {
	Status        string                 `json:"status"`
	InstanceID    string                 `json:"instanceId"`
	ProvisionedAt string                 `json:"provisionedAt"`
	Provisioned   bool                   `json:"provisioned"`
	Metadata      map[string]interface{} `json:"metadata"`
	System        map[string]interface{} `json:"system"`
}

func getBaseDir() string {
	if envDir := os.Getenv("YDB_SILLY_BASEDIR"); envDir != "" {
		return envDir
	}
	if _, err := os.Stat(filepath.Join(BaseDir, "md", "API_DOCS.md")); err == nil {
		return BaseDir
	}
	if _, err := os.Stat(filepath.Join("..", "md", "API_DOCS.md")); err == nil {
		return ".."
	}
	return BaseDir
}

func tryInitYDB() error {
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn != nil {
		return nil
	}
	var err error
	db, err = yottadb.Init()
	if err != nil {
		return fmt.Errorf("yottadb init failed: %w", err)
	}
	conn = yottadb.NewConn()
	m, err = conn.Import(mCallTable)
	if err != nil {
		return fmt.Errorf("yottadb import failed: %w", err)
	}
	return nil
}

func initYDB() {
	fmt.Println("--- ⚡ BOOTING YOTTADB NATIVE ENGINE (GO) ---")
	if err := tryInitYDB(); err != nil {
		fmt.Printf("[ENGINE] Initialization Failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[ENGINE] Connected.")
}

func getVMState(id string) *VMState {
	if id == "" || !IsDBConnected() {
		return nil
	}

	status := GetVM(id, "status")
	if status == "" {
		return nil
	}

	state := &VMState{
		Status:        status,
		InstanceID:    GetVM(id, "instanceId"),
		ProvisionedAt: GetVM(id, "provisionedAt"),
		Provisioned:   GetVM(id, "provisioned") == "1",
		Metadata: map[string]interface{}{
			"internalIp": GetMeta(id, "internalIp"),
			"port":       GetMeta(id, "port"),
			"region":     GetMeta(id, "region"),
			"cpuLimit":   GetMeta(id, "cpuLimit"),
		},
	}
	return state
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-VM-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// API Endpoints

	r.GET("/api/vm/status", func(c *gin.Context) {
		vmID := c.GetHeader("X-VM-ID")
		if vmID == "" {
			vmID = c.Query("id")
		}
		state := getVMState(vmID)
		if state == nil {
			state = &VMState{Status: "OFFLINE", Provisioned: false}
		}

		engineVer := yottadb.WrapperRelease
		if v, err := CallM("Version"); err == nil && v != nil {
			if s, ok := v.(string); ok && s != "" {
				engineVer = s
			}
		}

		state.System = map[string]interface{}{
			"engine":    engineVer,
			"nativeJob": os.Getpid(),
			"uptime":    GetUptime(),
		}

		c.JSON(http.StatusOK, state)
	})

	r.POST("/api/vm/provision", func(c *gin.Context) {
		var req struct {
			ID    string `json:"id"`
			Cores int    `json:"cores"`
			Sync  bool   `json:"sync"`
		}
		_ = c.ShouldBindJSON(&req)

		newID := strings.TrimSpace(req.ID)
		if newID == "" {
			newID = fmt.Sprintf("ydb-%d", rand.Intn(100000000))
		}
		cores := req.Cores
		if cores <= 0 {
			cores = []int{3, 5, 7}[rand.Intn(3)]
		}

		Provision(newID, fmt.Sprintf("%d", cores))
		LogInfo(fmt.Sprintf("Initiated provisioning for %s", newID))

		state := getVMState(newID)

		if req.Sync {
			Complete(newID)
			LogInfo(fmt.Sprintf("Provisioning complete for %s", newID))
			state = getVMState(newID)
		} else {
			go func(id string) {
				time.Sleep(2 * time.Second)
				Complete(id)
				LogInfo(fmt.Sprintf("Provisioning complete for %s", id))
			}(newID)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Provisioning initiated",
			"instanceId": newID,
			"state":      state,
		})
	})

	r.GET("/api/vm/list", func(c *gin.Context) {
		ids := ListVMs()
		vms := make([]*VMState, 0, len(ids))
		for _, id := range ids {
			if st := getVMState(id); st != nil {
				vms = append(vms, st)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"count": len(vms),
			"ids":   ids,
			"vms":   vms,
		})
	})

	r.DELETE("/api/vm/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "all" {
			CleanVM("")
		} else {
			CleanVM(id)
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "cleaned": id})
	})

	executeHandler := func(c *gin.Context) {
		var req struct {
			MCode string `json:"mCode"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.MCode) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		lineCount, err := XExe(req.MCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "FAIL"})
			return
		}

		var lines []string
		jobID := os.Getpid()

		for i := 1; i <= lineCount; i++ {
			val := GetGlobal("XOUT", jobID, i)
			lines = append(lines, val)
		}

		c.JSON(http.StatusOK, gin.H{
			"output": strings.TrimSpace(strings.Join(lines, "\n")),
			"lines":  lineCount,
			"status": "OK",
			"method": "NATIVE",
		})
	}

	r.POST("/execute", executeHandler)
	r.POST("/api/execute", executeHandler)

	r.GET("/api/files", func(c *gin.Context) {
		type FileInfo struct {
			Name   string `json:"name"`
			Folder string `json:"folder"`
			Path   string `json:"path"`
			Suffix string `json:"suffix"`
		}
		var allFiles []FileInfo
		base := getBaseDir()

		for _, folder := range folders {
			dirPath := filepath.Join(base, folder)
			entries, err := os.ReadDir(dirPath)
			if err == nil {
				for _, f := range entries {
					if !f.IsDir() {
						allFiles = append(allFiles, FileInfo{
							Name:   f.Name(),
							Folder: folder,
							Path:   filepath.Join(folder, f.Name()),
							Suffix: strings.ToLower(filepath.Ext(f.Name())),
						})
					}
				}
			}
		}

		// Root files
		entries, err := os.ReadDir(base)
		if err == nil {
			for _, f := range entries {
				if !f.IsDir() {
					allFiles = append(allFiles, FileInfo{
						Name:   f.Name(),
						Folder: "/",
						Path:   f.Name(),
						Suffix: strings.ToLower(filepath.Ext(f.Name())),
					})
				}
			}
		}

		// Sort by suffix ONLY (then name for deterministic tie-breaking)
		sort.Slice(allFiles, func(i, j int) bool {
			if allFiles[i].Suffix < allFiles[j].Suffix {
				return true
			}
			if allFiles[i].Suffix > allFiles[j].Suffix {
				return false
			}
			return allFiles[i].Name < allFiles[j].Name
		})

		c.JSON(http.StatusOK, gin.H{"files": allFiles})
	})

	r.GET("/api/global/:name", func(c *gin.Context) {
		name := normalizeGlobalName(c.Param("name"))
		subsStr := c.Query("subs")
		var subs []any
		if subsStr != "" {
			for _, s := range strings.Split(subsStr, ",") {
				subs = append(subs, s)
			}
		}

		val := GetGlobal(name, subs...)
		exists := GlobalExists(name, subs...)
		c.JSON(http.StatusOK, gin.H{
			"global":     name,
			"subscripts": subs,
			"value":      val,
			"exists":     exists,
		})
	})

	r.POST("/api/global/:name", func(c *gin.Context) {
		name := normalizeGlobalName(c.Param("name"))
		var req struct {
			Subscripts []string `json:"subscripts"`
			Value      string   `json:"value"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var subs []any
		for _, s := range req.Subscripts {
			subs = append(subs, s)
		}

		SetGlobal(name, req.Value, subs...)
		c.JSON(http.StatusOK, gin.H{"status": "OK", "global": name})
	})

	r.DELETE("/api/global/:name", func(c *gin.Context) {
		name := normalizeGlobalName(c.Param("name"))
		subsStr := c.Query("subs")
		var subs []any
		if subsStr != "" {
			for _, s := range strings.Split(subsStr, ",") {
				subs = append(subs, s)
			}
		}
		KillGlobal(name, subs...)
		c.JSON(http.StatusOK, gin.H{"status": "OK", "global": name})
	})

	r.GET("/api/global/:name/tree", func(c *gin.Context) {
		name := normalizeGlobalName(c.Param("name"))
		max := 50
		if mStr := c.Query("max"); mStr != "" {
			if parsed, err := strconv.Atoi(mStr); err == nil && parsed > 0 {
				max = parsed
			}
		}
		count := 0
		var entries []string
		if v, err := CallM("DataCount", name); err == nil && v != nil {
			if n, ok := v.(int); ok {
				count = n
			}
		}
		if v, err := CallM("DataDump", name, max); err == nil && v != nil {
			if s, ok := v.(string); ok && s != "" {
				entries = strings.Split(s, "\n")
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"global":  name,
			"count":   count,
			"entries": entries,
		})
	})

	r.GET("/api/logs", func(c *gin.Context) {
		lvl := strings.ToUpper(strings.TrimSpace(c.DefaultQuery("level", "INFO")))
		count := 10
		if cntStr := c.Query("count"); cntStr != "" {
			if parsed, err := strconv.Atoi(cntStr); err == nil && parsed > 0 {
				count = parsed
			}
		}
		raw := ReadLogs(lvl, count)
		var lines []string
		for _, line := range strings.Split(strings.TrimSpace(raw), "\n") {
			if line != "" {
				lines = append(lines, line)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"level": lvl,
			"logs":  lines,
			"stats": GetLogStats(),
		})
	})

	r.DELETE("/api/logs", func(c *gin.Context) {
		lvl := strings.ToUpper(strings.TrimSpace(c.Query("level")))
		ClearLogs(lvl)
		c.JSON(http.StatusOK, gin.H{"status": "OK", "cleared": lvl})
	})

	r.GET("/api/sys/intrinsics", func(c *gin.Context) {
		telemetry := map[string]string{
			"uptime":   GetUptime(),
			"datetime": GetDateTime(),
			"pid":      fmt.Sprintf("%d", os.Getpid()),
		}
		if v, err := CallM("Intrinsics"); err == nil && v != nil {
			if s, ok := v.(string); ok {
				for _, kv := range strings.Split(s, "|") {
					parts := strings.SplitN(kv, "=", 2)
					if len(parts) == 2 {
						telemetry[parts[0]] = parts[1]
					}
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":     "OK",
			"connected":  IsDBConnected(),
			"intrinsics": telemetry,
		})
	})

	r.GET("/api/cats", func(c *gin.Context) {
		if v, err := CallM("CatsJSON"); err == nil && v != nil {
			if s, ok := v.(string); ok && s != "" {
				c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(s))
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"cats": []gin.H{
				{"NAME": "MR. WHISKERS", "MOOD": "HISSING"},
				{"NAME": "GARFIELD", "MOOD": "LASAGNA"},
				{"NAME": "BINGUS", "MOOD": "BALD"},
			},
		})
	})

	r.POST("/api/cats", func(c *gin.Context) {
		var req struct {
			Name string `json:"name"`
			Mood string `json:"mood"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		id := 0
		if v, err := CallM("CatsAdd", req.Name, req.Mood); err == nil && v != nil {
			if n, ok := v.(int); ok {
				id = n
			}
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "id": id, "name": req.Name, "mood": req.Mood})
	})

	r.POST("/api/forge", func(c *gin.Context) {
		lines, err := CallM("ForgeEmoji")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "FAIL"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "FORGED",
			"routine": "ZEMOJI",
			"lines":   lines,
		})
	})

	r.GET("/api/decode/:char", func(c *gin.Context) {
		ch := c.Param("char")
		if decoded, err := URLDecode(ch); err == nil {
			ch = decoded
		}
		if v, err := CallM("DecodeEmoji", ch); err == nil && v != nil {
			if s, ok := v.(string); ok && s != "" {
				c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(s))
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"char": ch, "error": "ENGINE_OFFLINE"})
	})

	for _, sillyPath := range []string{"/api/wisdom", "/api/dogs", "/api/frogs"} {
		pathCopy := sillyPath
		r.GET(pathCopy, func(c *gin.Context) {
			if v, err := CallM("RigidRoute", pathCopy); err == nil && v != nil {
				if s, ok := v.(string); ok && s != "" {
					c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(s))
					return
				}
			}
			fallback := map[string]string{
				"/api/wisdom": "MUMPS IS ETERNAL",
				"/api/dogs":   "MAXIMUM BORK",
				"/api/frogs":  "REEEEEEEEEEEEE",
			}
			c.JSON(http.StatusOK, gin.H{"data": fallback[pathCopy]})
		})
	}

	r.GET("/api/mumps/test", func(c *gin.Context) {
		summary, err := CallM("RunMTests")
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "OFFLINE", "error": err.Error()})
			return
		}
		sumStr, _ := summary.(string)
		totalStr := GetGlobal("YDBTEST", "TOTAL")
		total, _ := strconv.Atoi(totalStr)
		var cases []string
		for i := 1; i <= total; i++ {
			if tc := GetGlobal("YDBTEST", "CASES", i); tc != "" {
				cases = append(cases, tc)
			}
		}
		status := "PASS"
		if !strings.Contains(sumStr, "FAIL:0") {
			status = "FAIL"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"summary": sumStr,
			"total":   total,
			"cases":   cases,
		})
	})

	// Documentation Route
	r.GET("/api/docs", func(c *gin.Context) {
		content, err := os.ReadFile(filepath.Join(getBaseDir(), "md", "API_DOCS.md"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Documentation not found")
			return
		}
		c.Data(http.StatusOK, "text/markdown; charset=utf-8", content)
	})

	r.GET("/v1beta/api", func(c *gin.Context) {
		content, _ := os.ReadFile(filepath.Join(getBaseDir(), "md", "API_DOCS.md"))
		c.String(http.StatusOK, string(content))
	})

	// Rigid MUMPS JSON healthcheck
	r.GET("/api/health", func(c *gin.Context) {
		health := "OK"
		if v, err := CallM("SysHealth"); err == nil && v != nil {
			if s, ok := v.(string); ok && s != "" {
				health = s
			}
		}
		raw, _ := json.Marshal(gin.H{
			"status":    health,
			"connected": IsDBConnected(),
			"uptime":    GetUptime(),
		})
		c.Data(http.StatusOK, "application/json; charset=utf-8", raw)
	})

	return r
}

func main() {
	initYDB()
	defer yottadb.Shutdown(db)

	r := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf(`
    =========================================
    🚀 YOTTADB NATIVE PROXY V4.0.0 (GO)
    PORT: %s
    MODE: NATIVE (GO-API)
    BONAFIDES: GOPHERIZED ✅
    =========================================
    `, port)

	r.Run(":" + port)
}
