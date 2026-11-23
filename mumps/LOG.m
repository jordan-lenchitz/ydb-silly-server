LOG ; 📝 NATIVE LOGGING ENGINE
 ; Persistent structured audit log in ^YDBLOGS
 ;
INFO(MSG) ; Log Info
 Q $$WRITE("INFO",$G(MSG))
 ;
WARN(MSG) ; Log Warning
 Q $$WRITE("WARN",$G(MSG))
 ;
ERROR(MSG) ; Log Error
 Q $$WRITE("ERROR",$G(MSG))
 ;
WRITE(LVL,MSG) ; Internal Atomic Writer
 N TS,ID
 S LVL=$G(LVL)
 I LVL="" S LVL="INFO"
 S TS=$$DATETIME^SYS()
 S ID=$I(^YDBLOGS(LVL))
 S ^YDBLOGS(LVL,ID,"ts")=TS
 S ^YDBLOGS(LVL,ID,"msg")=$G(MSG)
 S ^YDBLOGS(LVL,ID,"job")=$J
 S ^YDBLOGS(LVL,ID,"zut")=$ZUT
 Q ID
 ;
READ(LVL,COUNT) ; Read last N logs for level
 N ID,I,RES,ST,CNT
 S LVL=$G(LVL)
 I LVL="" Q ""
 S ID=+$G(^YDBLOGS(LVL))
 I ID<1 Q ""
 S CNT=$S(+$G(COUNT)>0:+COUNT,1:10)
 S ST=ID-CNT+1
 I ST<1 S ST=1
 S RES=""
 F I=ID:-1:ST D
 . S RES=RES_LVL_" ["_$G(^YDBLOGS(LVL,I,"ts"))_"] "_$G(^YDBLOGS(LVL,I,"msg"))_$C(10)
 Q RES
 ;
CLEAR(LVL) ; Clear logs
 I $G(LVL)="" K ^YDBLOGS Q "1"
 K ^YDBLOGS(LVL)
 Q "1"
 ;
STATS() ; Get log statistics across levels
 N L,C,RES
 S RES="",L=""
 F  S L=$O(^YDBLOGS(L)) Q:L=""  D
 . S C=+$G(^YDBLOGS(L))
 . S RES=RES_L_":"_C_";"
 Q RES
