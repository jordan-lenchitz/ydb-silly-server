DATA ; 💾 CORE PERSISTENCE LAYER
 ; Standardized access to YDB Globals
 ;
GET(GBL,SUBS) ; Generic Global Getter
 ; GBL: Global name (without ^)
 ; SUBS: Subscripts separated by |
 NEW S,CMD,VAL,I
 SET CMD="SET VAL=$G(^"_GBL_"("
 FOR I=1:1:$L(SUBS,"|") DO
 . SET S=$P(SUBS,"|",I)
 . IF S'="" SET CMD=CMD_""""_S_""""_($S(I<$L(SUBS,"|"):"," , 1:""))
 SET CMD=CMD_"))"
 XECUTE CMD
 QUIT VAL
 ;
SET(GBL,SUBS,VAL) ; Generic Global Setter
 NEW S,CMD,I
 SET CMD="SET ^"_GBL_"("
 FOR I=1:1:$L(SUBS,"|") DO
 . SET S=$P(SUBS,"|",I)
 . IF S'="" SET CMD=CMD_""""_S_""""_($S(I<$L(SUBS,"|"):"," , 1:""))
 SET CMD=CMD_")="""_VAL_""""
 XECUTE CMD
 QUIT 1
 ;
KILL(GBL,SUBS) ; Generic Global Killer
 NEW S,CMD,I
 SET CMD="KILL ^"_GBL_"("
 FOR I=1:1:$L(SUBS,"|") DO
 . SET S=$P(SUBS,"|",I)
 . IF S'="" SET CMD=CMD_""""_S_""""_($S(I<$L(SUBS,"|"):"," , 1:""))
 SET CMD=CMD_")"
 XECUTE CMD
 QUIT 1
 ;
EXISTS(GBL,SUBS) ; Check if node exists
 NEW S,CMD,RES,I
 SET CMD="SET RES=$D(^"_GBL_"("
 FOR I=1:1:$L(SUBS,"|") DO
 . SET S=$P(SUBS,"|",I)
 . IF S'="" SET CMD=CMD_""""_S_""""_($S(I<$L(SUBS,"|"):"," , 1:""))
 SET CMD=CMD_"))"
 XECUTE CMD
 QUIT RES
 ;
QUERY(GBL,SUBS) ; Simple query for next subscript
 NEW S,CMD,RES,I
 SET CMD="SET RES=$O(^"_GBL_"("
 FOR I=1:1:$L(SUBS,"|") DO
 . SET S=$P(SUBS,"|",I)
 . IF S'="" SET CMD=CMD_""""_S_""""_($S(I<$L(SUBS,"|"):"," , 1:""))
 SET CMD=CMD_"))"
 XECUTE CMD
 QUIT RES
