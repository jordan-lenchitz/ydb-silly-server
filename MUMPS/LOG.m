LOG ; 📝 LOGGING ENGINE
 ; Native logging system for YDB Cloud
 ;
INFO(MSG) ; Log Info
 DO WRITE("INFO",MSG)
 QUIT
 ;
WARN(MSG) ; Log Warning
 DO WRITE("WARN",MSG)
 QUIT
 ;
ERROR(MSG) ; Log Error
 DO WRITE("ERROR",MSG)
 QUIT
 ;
WRITE(LVL,MSG) ; Internal Writer
 NEW TS,ID
 SET TS=$H
 SET ID=$I(^YDBLOGS(LVL))
 SET ^YDBLOGS(LVL,ID,"ts")=TS
 SET ^YDBLOGS(LVL,ID,"msg")=MSG
 SET ^YDBLOGS(LVL,ID,"job")=$J
 QUIT
 ;
READ(LVL,COUNT) ; Read last N logs
 NEW ID,I,RES
 SET ID=$G(^YDBLOGS(LVL))
 IF ID="" QUIT ""
 SET RES=""
 FOR I=ID:-1:$S(ID-COUNT<1:1,1:ID-COUNT+1) DO
 . SET RES=RES_LVL_" ["_$G(^YDBLOGS(LVL,I,"ts"))_"] "_$G(^YDBLOGS(LVL,I,"msg"))_$C(10)
 QUIT RES
 ;
CLEAR(LVL) ; Clear logs
 IF LVL="" KILL ^YDBLOGS QUIT
 KILL ^YDBLOGS(LVL)
 QUIT
 ;
STATS() ; Get log statistics
 NEW L,C,RES
 SET RES="",L=""
 FOR  SET L=$O(^YDBLOGS(L)) QUIT:L=""  DO
 . SET C=$G(^YDBLOGS(L))
 . SET RES=RES_L_":"_C_";"
 QUIT RES
