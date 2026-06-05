package main

import (
	"fmt"
	"math/rand"
)

// --- VMMGR Utilities ---

func ProvisionVM(id string, cores string) string {
	ip := fmt.Sprintf("172.17.0.%d", rand.Intn(254)+1)
	region := "us-central1-gen2"
	cpuLimit := cores + ".0"

	SetGlobal("YDBCLOUD", "PROVISIONING", "VMS", id, "status")
	SetGlobal("YDBCLOUD", id, "VMS", id, "instanceId") // Wait, the MUMPS code was: ^YDBCLOUD("VMS",ID,ATTR)=VAL
	// Let me re-read VMMGR.m SETVM
	/*
	SETVM(ID,ATTR,VAL) ; Internal State Setter
	SET ^YDBCLOUD("VMS",ID,ATTR)=VAL
	QUIT
	;
	SETVM(ID,SUB,ATTR,VAL) ; Metadata Setter (Overloaded)
	SET ^YDBCLOUD("VMS",ID,SUB,ATTR)=VAL
	QUIT
	*/
	
	// Correcting subscripts
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

// Wait, I messed up the order of arguments in SetGlobal above.
// SetGlobal(gbl string, val string, subs ...any)

func ProvisionVMCorrected(id string, cores string) string {
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

// Let me just write the final version properly.

func Provision(id string, cores string) string {
	ip := fmt.Sprintf("172.17.0.%d", rand.Intn(254)+1)
	region := "us-central1-gen2"
	cpuLimit := cores + ".0"

	SetGlobal("YDBCLOUD", "PROVISIONING", id, "VMS", id, "status") // Wait, the MUMPS code: SET ^YDBCLOUD("VMS",ID,ATTR)=VAL
	// Global name is "YDBCLOUD"
	// Subscripts are "VMS", ID, ATTR
	
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
