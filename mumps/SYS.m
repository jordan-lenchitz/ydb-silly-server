SYS ; 🖥️ YOTTADB $Z* SYSTEM INTROSPECTION
 ; Accessing YottaDB engine and OS metrics via $Z* intrinsics
 ;
VERSION() ; Get YottaDB Version
 Q $ZV
 ;
PID() ; Get Process ID
 Q $J
 ;
GBLDIR() ; Get Global Directory Path
 Q $ZGBLDIR
 ;
ROUTINES() ; Get Routine Search Path
 Q $ZRO
 ;
CHSET() ; Get Character Set (UTF-8 / M)
 Q $ZCHSET
 ;
UPTIME() ; Seconds since midnight from $HOROLOG
 Q $P($H,",",2)
 ;
UTC() ; Universal time in microseconds since epoch
 Q $ZUT
 ;
HOROLOG() ; High-resolution horolog
 Q $ZH
 ;
ENV(VAR) ; Translate environment variable via $ZTRNLNM
 Q $ZTRNLNM($G(VAR))
 ;
DATETIME() ; Formatted date/time using YottaDB $ZDATE
 Q $ZD($H,"YYYY-MM-DD 24:60:SS")
 ;
DATE(D) ; Format Date component
 I $G(D)="" S D=$P($H,",",1)
 Q $ZD(D,"YYYY-MM-DD")
 ;
TIME(T) ; Format Time component
 I $G(T)="" S T=$P($H,",",2)
 Q $ZD("0,"_(+T),"24:60:SS")
 ;
HEALTH() ; System Health Check
 Q "OK"
 ;
METRICS() ; Packed telemetry string
 Q "zv="_$ZV_"|pid="_$J_"|chset="_$ZCHSET_"|gbldir="_$ZGBLDIR_"|zut="_$ZUT_"|zh="_$ZH_"|zl="_$ZL_"|pos="_$ZPOS
