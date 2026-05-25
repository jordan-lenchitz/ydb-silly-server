WEB ; 🌐 NATIVE WEB UTILITIES
 ; Helper functions for web-related data handling
 ;
URLENC(S) ; URL Encode
 NEW I,C,RES,HEX
 SET RES="",HEX="0123456789ABCDEF"
 FOR I=1:1:$L(S) DO
 . SET C=$E(S,I)
 . IF C?1(1A,1N,1".",1"-",1"_",1"~") SET RES=RES_C QUIT
 . SET RES=RES_"%"_$E(HEX,($A(C)\16+1))_$E(HEX,($A(C)#16+1))
 QUIT RES
 ;
URLDEC(S) ; URL Decode
 NEW I,C,RES
 SET RES=""
 FOR I=1:1:$L(S) DO
 . SET C=$E(S,I)
 . IF C="%" DO  QUIT
 . . SET RES=RES_$C($$HEXTODEC($E(S,I+1,I+2)))
 . . SET I=I+2
 . IF C="+" SET RES=RES_" " QUIT
 . SET RES=RES_C
 QUIT RES
 ;
HEXTODEC(H) ; Convert Hex to Dec
 NEW I,C,V,RES,HEX
 SET RES=0,HEX="0123456789ABCDEF",H=$$UP^STR(H)
 FOR I=1:1:$L(H) DO
 . SET C=$E(H,I)
 . SET V=$F(HEX,C)-2
 . SET RES=RES*16+V
 QUIT RES
 ;
PARSEHDR(H,OUT) ; Parse HTTP Headers
 ; H: Raw headers string, OUT: Output array
 NEW I,LINE,KEY,VAL
 KILL @OUT
 FOR I=1:1:$L(H,$C(10)) DO
 . SET LINE=$P(H,$C(10),I)
 . QUIT:LINE=""
 . SET KEY=$$TRIM^STR($P(LINE,":",1))
 . SET VAL=$$TRIM^STR($P(LINE,":",2,999))
 . IF KEY'="" SET @OUT@(KEY)=VAL
 QUIT
 ;
JSON(OBJ) ; Mock JSON Serializer (limited)
 NEW I,RES,KEY,VAL,FIRST
 SET RES="{",FIRST=1
 SET KEY=""
 FOR  SET KEY=$O(@OBJ@(KEY)) QUIT:KEY=""  DO
 . SET VAL=@OBJ@(KEY)
 . IF 'FIRST SET RES=RES_","
 . SET RES=RES_""""_KEY_""":"""_VAL_""""
 . SET FIRST=0
 SET RES=RES_"}"
 QUIT RES
