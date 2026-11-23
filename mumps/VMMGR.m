VMMGR ; 🚀 NATIVE VM ORCHESTRATOR
 ; Core logic for VM provisioning and lifecycle management in ^YDBCLOUD
 ;
PROVISION(ID,CORES) ; Initialize VM Provisioning
 N AT,IP,REGION,CPUL
 I $G(ID)="" S ID="ydb-"_($R(89999999)+10000000)
 I $G(CORES)="" S CORES=$P("3,5,7",",",$R(3)+1)
 S AT=$$DATETIME^SYS()
 S IP="172.17.0."_($R(254)+1)
 S REGION="us-central1-gen2"
 S CPUL=CORES_".0"
 ;
 S ^YDBCLOUD("VMS",ID,"status")="PROVISIONING"
 S ^YDBCLOUD("VMS",ID,"instanceId")=ID
 S ^YDBCLOUD("VMS",ID,"provisionedAt")=AT
 S ^YDBCLOUD("VMS",ID,"provisioned")="1"
 S ^YDBCLOUD("VMS",ID,"heartbeat")=$ZUT
 S ^YDBCLOUD("VMS",ID,"zh")=$ZH
 S ^YDBCLOUD("VMS",ID,"metadata","internalIp")=IP
 S ^YDBCLOUD("VMS",ID,"metadata","port")="8080"
 S ^YDBCLOUD("VMS",ID,"metadata","region")=REGION
 S ^YDBCLOUD("VMS",ID,"metadata","cpuLimit")=CPUL
 Q "SUCCESS"
 ;
COMPLETE(ID) ; Mark provisioning as complete
 I $G(ID)="" Q "0"
 S ^YDBCLOUD("VMS",ID,"status")="RUNNING"
 S ^YDBCLOUD("VMS",ID,"heartbeat")=$ZUT
 Q "1"
 ;
STOP(ID) ; Stop a running VM instance
 I $G(ID)="" Q "0"
 I '$D(^YDBCLOUD("VMS",ID)) Q "0"
 S ^YDBCLOUD("VMS",ID,"status")="STOPPED"
 S ^YDBCLOUD("VMS",ID,"heartbeat")=$ZUT
 Q "1"
 ;
HEARTBEAT(ID) ; Touch VM microsecond heartbeat via $ZUT
 I $G(ID)="" Q "0"
 S ^YDBCLOUD("VMS",ID,"heartbeat")=$ZUT
 Q $ZUT
 ;
SETVM(ID,ATTR,VAL) ; Internal State Setter
 I ($G(ID)="")!($G(ATTR)="") Q "0"
 S ^YDBCLOUD("VMS",ID,ATTR)=$G(VAL)
 Q "1"
 ;
SETMETA(ID,ATTR,VAL) ; Metadata Setter
 I ($G(ID)="")!($G(ATTR)="") Q "0"
 S ^YDBCLOUD("VMS",ID,"metadata",ATTR)=$G(VAL)
 Q "1"
 ;
GETVM(ID,ATTR) ; Get VM Attribute
 I ($G(ID)="")!($G(ATTR)="") Q ""
 Q $G(^YDBCLOUD("VMS",ID,ATTR))
 ;
GETMETA(ID,ATTR) ; Get VM Metadata
 I ($G(ID)="")!($G(ATTR)="") Q ""
 Q $G(^YDBCLOUD("VMS",ID,"metadata",ATTR))
 ;
LIST() ; Return comma-separated list of VM IDs
 N ID,RES S ID="",RES=""
 F  S ID=$O(^YDBCLOUD("VMS",ID)) Q:ID=""  D
 . S RES=RES_$S(RES="":"",1:",")_ID
 Q RES
 ;
COUNT() ; Return total number of VMs in ^YDBCLOUD
 N ID,C S ID="",C=0
 F  S ID=$O(^YDBCLOUD("VMS",ID)) Q:ID=""  S C=C+1
 Q C
 ;
CLEAN(ID) ; Remove VM record(s)
 I $G(ID)="" K ^YDBCLOUD("VMS") Q "1"
 K ^YDBCLOUD("VMS",ID)
 Q "1"
