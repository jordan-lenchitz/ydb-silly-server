WEB ; 🌐 NATIVE WEB, URL & JSON UTILITIES
 ; Byte-perfect URL encoding/decoding using $ZLENGTH, $ZASCII, $ZCHAR
 ;
URLENC(S) ; UTF-8 Byte-Perfect URL Encode
 N I,B,C,RES,HEX
 S S=$G(S),RES="",HEX="0123456789ABCDEF"
 F I=1:1:$ZLENGTH(S) D
 . S B=$ZASCII(S,I),C=$ZEXTRACT(S,I)
 . I (B'<48&(B'>57))!(B'<65&(B'>90))!(B'<97&(B'>122))!(C=".")!(C="-")!(C="_")!(C="~") S RES=RES_C Q
 . I C=" " S RES=RES_"+" Q
 . S RES=RES_"%"_$E(HEX,(B\16)+1)_$E(HEX,(B#16)+1)
 Q RES
 ;
URLDEC(S) ; UTF-8 Byte-Perfect URL Decode
 N I,C,RES,H1,H2
 S S=$G(S),RES=""
 F I=1:1:$ZLENGTH(S) D
 . S C=$ZEXTRACT(S,I)
 . I C="%" D  Q
 . . S H1=$ZEXTRACT(S,I+1,I+2)
 . . I $ZLENGTH(H1)=2 S RES=RES_$ZCHAR($$HEXTODEC(H1)),I=I+2 Q
 . . S RES=RES_C
 . I C="+" S RES=RES_" " Q
 . S RES=RES_C
 Q RES
 ;
HEXTODEC(H) ; Convert Hex to Decimal
 Q $$HD^MATH($G(H))
 ;
DECTOHEX(D) ; Convert Decimal to Hex
 Q $$DH^MATH(+$G(D))
 ;
PARSEHDR(H,OUT) ; Parse HTTP Headers into @OUT@(KEY)=VAL
 N I,LINE,KEY,VAL,CNT
 S H=$G(H),CNT=0
 K @OUT
 F I=1:1:$L(H,$C(10)) D
 . S LINE=$$TRIM^STR($P(H,$C(10),I))
 . Q:LINE=""
 . Q:LINE'[":"
 . S KEY=$$TRIM^STR($P(LINE,":",1))
 . S VAL=$$TRIM^STR($P(LINE,":",2,999))
 . I KEY'="" S @OUT@(KEY)=VAL,CNT=CNT+1
 Q:$Q CNT
 Q
 ;
HDRGET(H,TARGET) ; Extract a single header value (case-insensitive)
 N HMAP,K,T
 D PARSEHDR($G(H),"HMAP")
 S T=$$UP^STR($$TRIM^STR($G(TARGET))),K=""
 F  S K=$O(HMAP(K)) Q:K=""  I $$UP^STR(K)=T Q
 Q $S(K="":"",1:$G(HMAP(K)))
 ;
JSONESC(S) ; Escape string for JSON
 N I,C,R S S=$G(S),R=""
 F I=1:1:$L(S) D
 . S C=$E(S,I)
 . I C="""" S R=R_"\""" Q
 . I C="\" S R=R_"\\" Q
 . I C=$C(10) S R=R_"\n" Q
 . I C=$C(13) S R=R_"\r" Q
 . I C=$C(9) S R=R_"\t" Q
 . S R=R_C
 Q R
 ;
JSON(OBJ) ; Serialize flat subscripted array @OBJ to JSON object
 N KEY,VAL,RES,FIRST
 S RES="{",FIRST=1,KEY=""
 F  S KEY=$O(@OBJ@(KEY)) Q:KEY=""  D
 . S VAL=$G(@OBJ@(KEY))
 . I 'FIRST S RES=RES_","
 . S RES=RES_""""_$$JSONESC(KEY)_""":"""_$$JSONESC(VAL)_""""
 . S FIRST=0
 Q RES_"}"
