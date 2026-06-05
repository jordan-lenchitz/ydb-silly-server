package main

import (
	"fmt"
	"math/rand"
)

// --- VMMGR Utilities ---

func Provision(id string, cores string) string {
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
	SetGlobal("YDBCLOUD", "RUNNING", "VMS", id, "status")
	return "1"
}

func GetVM(id, attr string) string {
	return GetGlobal("YDBCLOUD", "VMS", id, attr)
}

func GetMeta(id, attr string) string {
	return GetGlobal("YDBCLOUD", "VMS", id, "metadata", attr)
}

func CleanVM(id string) {
	if id == "" {
		KillGlobal("YDBCLOUD", "VMS")
	} else {
		KillGlobal("YDBCLOUD", "VMS", id)
	}
}
