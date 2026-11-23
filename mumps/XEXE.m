BOOT ; 🚀 CLI ENTRY POINT FOR ydb -run ^XEXE
 N %V,%N
 S %V=$$VERSION(),%N=$$XEXE("W ""[XEXE] SELF-TEST OK FROM JOB "",$J,!")
 W "--- ⚡ YDB-SILLY-SERVER MUMPS ENGINE ---",!
 W "VERSION:    ",%V,!
 W "INTRINSICS: ",$$INTRINSICS(),!
 W "LINES:      ",%N," -> ",$G(^XOUT($J,1)),!
 Q
 ;
XEXE(CODE) ; 🚀 THE UNHOLY EXECUTIONER
 N %FN,%LINE,%OUT,%I
 S %OUT=$J,%I=0
 K ^XOUT(%OUT)
 S %FN="/tmp/ydb_out_"_%OUT_"_"_$R(999999)
 O %FN:(newversion:truncate:stream):0
 D RUN(CODE,%FN)
 O %FN:(readonly:stream):0
 F  U %FN R %LINE D:$L(%LINE)!('$ZEOF)  Q:$ZEOF
 . S %I=%I+1,^XOUT(%OUT,%I)=%LINE
 C %FN:(DELETE)
 U $P
 Q:$Q %I
 Q
 ;
RUN(CODE,%FN) ; Isolated stack frame with $ETRAP error capture
 N $ET S $ET="G XERR^XEXE"
 U %FN
 X CODE
 C %FN
 Q
 ;
XERR ; Error handler for unsafe execution
 U %FN W "[MUMPS-ERROR] "_$ZS,! C %FN
 S $EC=""
 Q
 ;
EVAL(EXPR) ; Evaluate arbitrary MUMPS expression and return string
 N %R S %R=""
 D EVRUN($G(EXPR),.%R)
 Q %R
 ;
EVRUN(EXPR,%R) ; Isolated stack frame for expression evaluation
 N $ET S $ET="G EVERR^XEXE"
 X "S %R="_EXPR
 Q
 ;
EVERR ; Expression evaluation error trap
 S %R="[EVAL-ERROR] "_$ZS,$EC=""
 Q
 ;
VERSION() ; Get YottaDB release & GT.M version
 Q $ZYRELEASE_" ("_$ZV_")"
 ;
INTRINSICS() ; Dump YottaDB-specific $Z* intrinsic telemetry
 Q "zyrelease="_$ZYRELEASE_"|zv="_$ZV_"|zchset="_$ZCHSET_"|zgbldir="_$ZGBLDIR_"|zro="_$ZRO_"|zh="_$ZH_"|zut="_$ZUT_"|zpos="_$ZPOS_"|zl="_$ZL_"|job="_$J
