MATH ; 🔢 NATIVE MATHEMATICS LIBRARY
 ; Optimized math routines for YottaDB
 ;
FACT(N) ; Factorial
 NEW I,RES
 SET RES=1
 FOR I=1:1:N SET RES=RES*I
 QUIT RES
 ;
ISPRIME(N) ; Check if number is prime
 NEW I,RES
 IF N<2 QUIT 0
 IF N=2 QUIT 1
 IF N#2=0 QUIT 0
 SET RES=1
 FOR I=3:2:($SQRT(N)\1) IF N#I=0 SET RES=0 QUIT
 QUIT RES
 ;
GCD(A,B) ; Greatest Common Divisor
 NEW TEMP
 FOR  QUIT:B=0  SET TEMP=B,B=A#B,A=TEMP
 QUIT A
 ;
LCM(A,B) ; Least Common Multiple
 QUIT (A*B)/$$GCD(A,B)
 ;
ABS(N) ; Absolute Value
 QUIT $S(N<0:-N,1:N)
 ;
POW(B,E) ; Power
 NEW I,RES
 SET RES=1
 FOR I=1:1:E SET RES=RES*B
 QUIT RES
 ;
ROUND(N,D) ; Round to D decimal places
 NEW P
 SET P=$$POW(10,D)
 QUIT (N*P+0.5\1)/P
 ;
FIB(N) ; Fibonacci
 IF N<2 QUIT N
 QUIT $$FIB(N-1)+$$FIB(N-2)
 ;
AVG(LIST) ; Average of a subscripted array
 NEW I,S,C
 SET (I,S,C)=0
 FOR  SET I=$O(@LIST@(I)) QUIT:I=""  SET S=S+@LIST@(I),C=C+1
 QUIT $S(C=0:0,1:S/C)
