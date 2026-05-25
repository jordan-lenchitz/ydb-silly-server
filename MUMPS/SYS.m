SYS ; 🖥️ SYSTEM INFORMATION UTILITIES
 ; Accessing YottaDB and OS metrics
 ;
VERSION() ; Get YDB Version
 QUIT $ZV
 ;
PID() ; Get Process ID
 QUIT $J
 ;
HOME() ; Get YDB Dist
 QUIT $G(^%ZSYS("DIST"))
 ;
GBLDIR() ; Get Global Directory
 QUIT $ZGBLDIR
 ;
UPTIME() ; Mock Uptime (seconds since midnight)
 QUIT $P($H,",",2)
 ;
MEMORY() ; Mock Memory Usage
 QUIT "512MB"
 ;
HOST() ; Get Hostname
 NEW H
 SET H=$S($D(^%ZSYS("HOST")):^%ZSYS("HOST"),1:"unknown")
 QUIT H
 ;
ENV(VAR) ; Get Environment Variable
 QUIT $ZTRNLNM(VAR)
 ;
DATETIME() ; Get formatted date/time
 NEW H,D,T
 SET H=$H,D=$P(H,",",1),T=$P(H,",",2)
 QUIT $$DATE(D)_" "_$$TIME(T)
 ;
DATE(D) ; Format Date
 QUIT $ZD(D,"YYYY-MM-DD")
 ;
TIME(T) ; Format Time
 QUIT $ZT(T)
 ;
HEALTH() ; System Health Check
 QUIT "OK"
