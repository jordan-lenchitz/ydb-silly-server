STR ; 🧵 UTF-8 & BYTE-PERFECT STRING UTILITIES
 ; High-performance string operations using $ZCONVERT, $ZLENGTH, $ZEXTRACT
 ;
REV(S) ; Reverse a string (UTF-8 character aware)
 N O,I S O="",S=$G(S)
 F I=$L(S):-1:1 S O=O_$E(S,I)
 Q O
 ;
UP(S) ; Uppercase via YottaDB $ZCONVERT
 Q $ZCONVERT($G(S),"U")
 ;
LOW(S) ; Lowercase via YottaDB $ZCONVERT
 Q $ZCONVERT($G(S),"L")
 ;
TRIM(S) ; Trim leading and trailing whitespace
 N F,L,C S S=$G(S)
 I S="" Q ""
 F F=1:1:$L(S) S C=$E(S,F) Q:(C'=" ")&(C'=$C(9))&(C'=$C(10))&(C'=$C(13))
 I F=$L(S) S C=$E(S,F) I (C=" ")!(C=$C(9))!(C=$C(10))!(C=$C(13)) Q ""
 F L=$L(S):-1:F S C=$E(S,L) Q:(C'=" ")&(C'=$C(9))&(C'=$C(10))&(C'=$C(13))
 Q $E(S,F,L)
 ;
SPLIT(S,D,OUT) ; Split string S by delimiter D into array @OUT
 N I,V S S=$G(S),D=$G(D,",")
 K @OUT
 I S="" Q 0
 F I=1:1:$L(S,D) S V=$P(S,D,I),@OUT@(I)=V
 Q I
 ;
JOIN(IN,D) ; Join subscripted array @IN with delimiter D
 N I,O,V,FIRST S O="",I="",D=$G(D,","),FIRST=1
 F  S I=$O(@IN@(I)) Q:I=""  D
 . S V=@IN@(I)
 . S O=O_$S(FIRST:"",1:D)_V,FIRST=0
 Q O
 ;
REPLACE(S,O,N) ; Replace all occurrences of O with N in S
 S S=$G(S),O=$G(O),N=$G(N)
 I O="" Q S
 N I,R S R=""
 F I=1:1:$L(S,O) S R=R_$P(S,O,I)_$S(I<$L(S,O):N,1:"")
 Q R
 ;
LEN(S) ; Character length
 Q $L($G(S))
 ;
BLEN(S) ; Raw UTF-8 byte length via $ZLENGTH
 Q $ZLENGTH($G(S))
 ;
SUB(S,ST,EN) ; Safe 1-based character substring
 S S=$G(S),ST=+$G(ST),EN=+$G(EN)
 I ST<1 S ST=1
 I EN>$L(S) S EN=$L(S)
 I ST>EN Q ""
 Q $E(S,ST,EN)
 ;
BSUB(S,ST,EN) ; Raw byte-level substring via $ZEXTRACT
 S S=$G(S),ST=+$G(ST),EN=+$G(EN)
 I ST<1 S ST=1
 I EN>$ZLENGTH(S) S EN=$ZLENGTH(S)
 I ST>EN Q ""
 Q $ZEXTRACT(S,ST,EN)
