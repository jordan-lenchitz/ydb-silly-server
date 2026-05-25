STR ; 🧵 STRING MANIPULATION UTILITIES
 ; High-performance string operations for M
 ;
REV(S) ; Reverse a string
 NEW O,I
 SET O=""
 FOR I=$L(S):-1:1 SET O=O_$E(S,I)
 QUIT O
 ;
UP(S) ; Uppercase
 QUIT $TR(S,"abcdefghijklmnopqrstuvwxyz","ABCDEFGHIJKLMNOPQRSTUVWXYZ")
 ;
LOW(S) ; Lowercase
 QUIT $TR(S,"ABCDEFGHIJKLMNOPQRSTUVWXYZ","abcdefghijklmnopqrstuvwxyz")
 ;
TRIM(S) ; Trim whitespace
 NEW F,L
 FOR F=1:1:$L(S) QUIT:$E(S,F)'=" "
 FOR L=$L(S):-1:1 QUIT:$E(S,L)'=" "
 QUIT $E(S,F,L)
 ;
SPLIT(S,D,OUT) ; Split string into array
 ; S: String, D: Delimiter, OUT: Output global/array name
 NEW I,V
 KILL @OUT
 FOR I=1:1:$L(S,D) DO
 . SET V=$P(S,D,I)
 . SET @OUT@(I)=V
 QUIT I
 ;
JOIN(IN,D) ; Join array into string
 NEW I,O,V
 SET O="",I=""
 FOR  SET I=$O(@IN@(I)) QUIT:I=""  DO
 . SET V=@IN@(I)
 . SET O=O_$S(O="":"",1:D)_V
 QUIT O
 ;
REPLACE(S,O,N) ; Replace Old with New
 NEW I,RES
 SET RES=""
 FOR I=1:1:$L(S,O) DO
 . SET RES=RES_$P(S,O,I)_$S(I<$L(S,O):N,1:"")
 QUIT RES
 ;
LEN(S) ; Length
 QUIT $L(S)
 ;
SUB(S,ST,EN) ; Substring
 QUIT $E(S,ST,EN)
