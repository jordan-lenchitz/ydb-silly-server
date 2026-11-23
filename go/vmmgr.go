package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// --- VMMGR Utilities ---

func Provision(id string, cores string) string {
	if id == "" {
		id = fmt.Sprintf("ydb-%d", rand.Intn(100000000))
	}
	if cores == "" {
		cores = "3"
	}
	if v, err := CallM("VMProvision", id, cores); err == nil && v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}

	ip := fmt.Sprintf("172.17.0.%d", rand.Intn(254)+1)
	region := "us-central1-gen2"
	cpuLimit := cores + ".0"

	SetGlobal("YDBCLOUD", "PROVISIONING", "VMS", id, "status")
	SetGlobal("YDBCLOUD", id, "VMS", id, "instanceId")
	SetGlobal("YDBCLOUD", GetDateTime(), "VMS", id, "provisionedAt")
	SetGlobal("YDBCLOUD", ip, "VMS", id, "metadata", "internalIp")
	SetGlobal("YDBCLOUD", "8080", "VMS", id, "metadata", "port")
	SetGlobal("YDBCLOUD", region, "VMS", id, "metadata", "region")
	SetGlobal("YDBCLOUD", cpuLimit, "VMS", id, "metadata", "cpuLimit")
	SetGlobal("YDBCLOUD", "1", "VMS", id, "provisioned")

	return "SUCCESS"
}

func Complete(id string) string {
	if v, err := CallM("VMComplete", id); err == nil && v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	SetGlobal("YDBCLOUD", "RUNNING", "VMS", id, "status")
	return "1"
}

func StopVM(id string) string {
	if v, err := CallM("VMStop", id); err == nil && v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	SetGlobal("YDBCLOUD", "STOPPED", "VMS", id, "status")
	return "1"
}

func GetVM(id, attr string) string {
	if v, err := CallM("VMGet", id, attr); err == nil && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return GetGlobal("YDBCLOUD", "VMS", id, attr)
}

func GetMeta(id, attr string) string {
	if v, err := CallM("VMMeta", id, attr); err == nil && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return GetGlobal("YDBCLOUD", "VMS", id, "metadata", attr)
}

func ListVMs() []string {
	if v, err := CallM("VMList"); err == nil && v != nil {
		if s, ok := v.(string); ok {
			if s == "" {
				return nil
			}
			return strings.Split(s, ",")
		}
	}
	dbMu.Lock()
	defer dbMu.Unlock()
	if conn == nil {
		return nil
	}
	var ids []string
	for _, id := range conn.Node("^YDBCLOUD", "VMS").Children() {
		ids = append(ids, id)
	}
	return ids
}

func CleanVM(id string) {
	if _, err := CallM("VMClean", id); err == nil {
		return
	}
	if id == "" {
		KillGlobal("YDBCLOUD", "VMS")
	} else {
		KillGlobal("YDBCLOUD", "VMS", id)
	}
}
