package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"lang.yottadb.com/go/yottadb/v2"
)

/**
 * 🚀 YOTTADB UNIRONIC BACKEND V4.0.0 (GO NATIVE)
 * Pure Gopher logic, YottaDB Persistence.
 */

const (
	BaseDir = "/app"
)

var (
	folders = []string{"MUMPS", "JS", "MD", "logs"}
	db      *yottadb.DB
	conn    *yottadb.Conn
	m       *yottadb.MFunctions
)

// Minimal call table for remaining M logic (XExe)
const mCallTable = `
	XExe: int XEXE^XEXE(string)
	Version: string VERSION^XEXE()
`

type VMState struct {
	Status        string                 `json:"status"`
	InstanceID    string                 `json:"instanceId"`
	ProvisionedAt string                 `json:"provisionedAt"`
	Provisioned   bool                   `json:"provisioned"`
	Metadata      map[string]interface{} `json:"metadata"`
	System        map[string]interface{} `json:"system"`
}

func initYDB() {
	fmt.Println("--- ⚡ BOOTING YOTTADB NATIVE ENGINE (GO) ---")
	var err error
	db, err = yottadb.Init()
	if err != nil {
		fmt.Printf("[ENGINE] Initialization Failed: %v\n", err)
		os.Exit(1)
	}
	conn = yottadb.NewConn()
	m, err = conn.Import(mCallTable)
	if err != nil {
		fmt.Printf("[ENGINE] Import Failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("[ENGINE] Connected.")
}

func getVMState(id string) *VMState {
	if id == "" || conn == nil {
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
		state := getVMState(vmID)
		if state == nil {
			state = &VMState{Status: "OFFLINE", Provisioned: false}
		}

		engineVer := yottadb.WrapperRelease
		if m != nil {
			if v, err := m.CallErr("Version"); err == nil && v != nil {
				engineVer = v.(string)
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
		newID := fmt.Sprintf("ydb-%d", rand.Intn(100000000))
		cores := []int{3, 5, 7}[rand.Intn(3)]

		Provision(newID, fmt.Sprintf("%d", cores))
		LogInfo(fmt.Sprintf("Initiated provisioning for %s", newID))

		state := getVMState(newID)

		go func() {
			time.Sleep(2 * time.Second)
			Complete(newID)
			LogInfo(fmt.Sprintf("Provisioning complete for %s", newID))
		}()

		c.JSON(http.StatusOK, gin.H{
			"message":    "Provisioning initiated",
			"instanceId": newID,
			"state":      state,
		})
	})

	executeHandler := func(c *gin.Context) {
		var req struct {
			MCode string `json:"mCode"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		lineCount, err := XExe(req.MCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "FAIL"})
			return
		}

		output := ""
		jobID := os.Getpid()

		for i := 1; i <= lineCount; i++ {
			val := GetGlobal("XOUT", jobID, i)
			output += val + "\n"
		}

		c.JSON(http.StatusOK, gin.H{
			"output": strings.TrimSpace(output),
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

		for _, folder := range folders {
			dirPath := filepath.Join(BaseDir, folder)
			files, err := ioutil.ReadDir(dirPath)
			if err == nil {
				for _, f := range files {
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
		files, err := ioutil.ReadDir(BaseDir)
		if err == nil {
			for _, f := range files {
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

		// Sort by suffix ONLY
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
		name := c.Param("name")
		subsStr := c.Query("subs")
		var subs []any
		if subsStr != "" {
			for _, s := range strings.Split(subsStr, ",") {
				subs = append(subs, s)
			}
		}

		val := GetGlobal(name, subs...)
		c.JSON(http.StatusOK, gin.H{"global": name, "subscripts": subs, "value": val})
	})

	r.POST("/api/global/:name", func(c *gin.Context) {
		name := c.Param("name")
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
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Documentation Route
	r.GET("/api/docs", func(c *gin.Context) {
		content, err := ioutil.ReadFile(filepath.Join(BaseDir, "MD/API_DOCS.md"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Documentation not found")
			return
		}
		c.Data(http.StatusOK, "text/markdown; charset=utf-8", content)
	})

	r.GET("/v1beta/api", func(c *gin.Context) {
		content, _ := ioutil.ReadFile(filepath.Join(BaseDir, "MD/API_DOCS.md"))
		c.String(http.StatusOK, string(content))
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
