VMMGR ; 🚀 NATIVE VM ORCHESTRATOR
 ; Core logic for VM provisioning and management
 ;
PROVISION(ID,CORES) ; Initialize VM Provisioning
 ; ID: Unique VM Identifier
 ; CORES: CPU Allocation
 NEW AT,IP,REGION,CPUL,STATUS
 SET AT=$H
 SET IP="172.17.0."_($R(254)+1)
 SET REGION="us-central1-gen2"
 SET CPUL=CORES_".0"
 ;
 ; Atomic state initialization
 DO SETVM(ID,"status","PROVISIONING")
 DO SETVM(ID,"instanceId",ID)
 DO SETVM(ID,"provisionedAt",$P(AT,",",1)_","_$P(AT,",",2))
 DO SETVM(ID,"metadata","internalIp",IP)
 DO SETVM(ID,"metadata","port",8080)
 DO SETVM(ID,"metadata","region",REGION)
 DO SETVM(ID,"metadata","cpuLimit",CPUL)
 DO SETVM(ID,"provisioned",1)
 ;
 QUIT "SUCCESS"
 ;
COMPLETE(ID) ; Mark provisioning as complete
 DO SETVM(ID,"status","RUNNING")
 QUIT 1
 ;
SETVM(ID,ATTR,VAL) ; Internal State Setter
 SET ^YDBCLOUD("VMS",ID,ATTR)=VAL
 QUIT
 ;
SETVM(ID,SUB,ATTR,VAL) ; Metadata Setter (Overloaded)
 SET ^YDBCLOUD("VMS",ID,SUB,ATTR)=VAL
 QUIT
 ;
GETVM(ID,ATTR) ; Get VM Attribute
 QUIT $G(^YDBCLOUD("VMS",ID,ATTR))
 ;
GETMETA(ID,ATTR) ; Get VM Metadata
 QUIT $G(^YDBCLOUD("VMS",ID,"metadata",ATTR))
 ;
LIST ; List all VMs
 NEW ID SET ID=""
 FOR  SET ID=$O(^YDBCLOUD("VMS",ID)) QUIT:ID=""  WRITE ID,!
 QUIT
 ;
CLEAN(ID) ; Remove VM record
 IF ID="" KILL ^YDBCLOUD("VMS") QUIT
 KILL ^YDBCLOUD("VMS",ID)
 QUIT
