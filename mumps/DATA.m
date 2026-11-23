DATA ; 💾 CORE B*-TREE PERSISTENCE LAYER
 ; Standardized access to YottaDB Globals with indirect references
 ;
CLEANGBL(GBL) ; Strip leading ^ and validate global name
 S GBL=$G(GBL)
 I $E(GBL,1)="^" S GBL=$E(GBL,2,$L(GBL))
 Q GBL
 ;
QUOTE(S) ; Double-quote escape a subscript string
 N I,C,R S R=""""
 F I=1:1:$L(S) S C=$E(S,I),R=R_$S(C="""":"""""",1:C)
 Q R_""""
 ;
REF(GBL,SUBS) ; Build indirect global reference ^GBL("s1","s2")
 N G,I,S,R,SEP
 S G=$$CLEANGBL($G(GBL))
 I G="" Q ""
 S SUBS=$G(SUBS)
 I SUBS="" Q "^"_G
 S R="^"_G_"(",SEP=""
 F I=1:1:$L(SUBS,"|") D
 . S S=$P(SUBS,"|",I)
 . S R=R_SEP_$$QUOTE(S),SEP=","
 Q R_")"
 ;
GET(GBL,SUBS) ; Generic Global Getter
 N R S R=$$REF($G(GBL),$G(SUBS))
 I R="" Q ""
 Q $G(@R)
 ;
SET(GBL,SUBS,VAL) ; Generic Global Setter
 N R S R=$$REF($G(GBL),$G(SUBS))
 I R="" Q 0
 S @R=$G(VAL)
 Q 1
 ;
KILL(GBL,SUBS) ; Generic Global Killer
 N R S R=$$REF($G(GBL),$G(SUBS))
 I R="" Q 0
 K @R
 Q 1
 ;
EXISTS(GBL,SUBS) ; Return $DATA value (0, 1, 10, 11)
 N R S R=$$REF($G(GBL),$G(SUBS))
 I R="" Q 0
 Q $D(@R)
 ;
INCR(GBL,SUBS,BY) ; Atomic $INCREMENT on a global node
 N R,STEP S R=$$REF($G(GBL),$G(SUBS))
 I R="" Q 0
 S STEP=$S(+$G(BY)=0:1,1:+BY)
 Q $I(@R,STEP)
 ;
ORDER(GBL,SUBS,DIR) ; Next/previous sibling subscript via $ORDER
 N R,D S R=$$REF($G(GBL),$G(SUBS))
 I R="" Q ""
 S D=$S(+$G(DIR)=-1:-1,1:1)
 Q $O(@R,D)
 ;
COUNT(GBL) ; Count all value-bearing nodes in global B*-tree via $QUERY
 N G,R,N,CNT
 S G=$$CLEANGBL($G(GBL))
 I G="" Q 0
 S R="^"_G,CNT=0
 I $D(@R)#2 S CNT=CNT+1
 S N=R
 F  S N=$Q(@N) Q:N=""  S CNT=CNT+1
 Q CNT
 ;
DUMP(GBL,MAX) ; Walk B*-tree via $QUERY and return newline-separated k=v pairs
 N G,R,N,CNT,LIM,OUT
 S G=$$CLEANGBL($G(GBL))
 I G="" Q ""
 S LIM=$S(+$G(MAX)>0:+MAX,1:50),CNT=0,OUT=""
 S R="^"_G
 I $D(@R)#2 S CNT=CNT+1,OUT=R_"="_@R
 S N=R
 F  S N=$Q(@N) Q:(N="")!(CNT'<LIM)  D
 . S CNT=CNT+1
 . S OUT=OUT_$S(OUT="":"",1:$C(10))_N_"="_@N
 Q OUT
