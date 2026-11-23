MATH ; 🔢 NATIVE MATHEMATICS & $ZBIT* BITWISE LIBRARY
 ; Optimized math and YottaDB intrinsic bitstring routines
 ;
FACT(N) ; Factorial
 N I,R S R=1
 F I=1:1:+$G(N) S R=R*I
 Q R
 ;
ISPRIME(N) ; Primality test
 N I,R S N=+$G(N)
 I N<2 Q 0
 I N=2 Q 1
 I N#2=0 Q 0
 S R=1
 F I=3:2 Q:(I*I)>N  I N#I=0 S R=0 Q
 Q R
 ;
GCD(A,B) ; Greatest Common Divisor
 N T S A=$$ABS(+$G(A)),B=$$ABS(+$G(B))
 F  Q:B=0  S T=B,B=A#B,A=T
 Q A
 ;
LCM(A,B) ; Least Common Multiple
 S A=+$G(A),B=+$G(B)
 I (A=0)!(B=0) Q 0
 Q $$ABS(A*B)\$$GCD(A,B)
 ;
ABS(N) ; Absolute Value
 S N=+$G(N)
 Q $S(N<0:-N,1:N)
 ;
POW(B,E) ; Integer Power
 N I,R S B=+$G(B),E=+$G(E),R=1
 I E<0 Q 0
 F I=1:1:E S R=R*B
 Q R
 ;
ROUND(N,D) ; Round N to D decimal places
 Q +$J(+$G(N),0,+$G(D))
 ;
FIB(N) ; Fibonacci with ^YMATH("FIB") global memoization
 S N=+$G(N)
 I N'>0 Q 0
 I N=1 Q 1
 I $D(^YMATH("FIB",N)) Q ^YMATH("FIB",N)
 N I,A,B,C S A=0,B=1
 F I=2:1:N S C=A+B,A=B,B=C,^YMATH("FIB",I)=C
 Q B
 ;
DH(D) ; Decimal to Hexadecimal
 N X,R,Y S X="0123456789ABCDEF",R="",D=+$G(D)
 I D'>0 Q "0"
 F  Q:D=0  S Y=D#16,R=$E(X,Y+1)_R,D=D\16
 Q R
 ;
HD(H) ; Hexadecimal to Decimal
 N I,C,V,R,X S R=0,X="0123456789ABCDEF",H=$ZCONVERT($G(H),"U")
 F I=1:1:$L(H) D
 . S C=$E(H,I),V=$F(X,C)-2
 . I V'<0 S R=(R*16)+V
 Q R
 ;
I2BS(N) ; Convert 16-bit non-negative integer to YottaDB $ZBITSTR
 N BS,I,V S V=+$G(N),BS=$ZBITSTR(16,0)
 F I=1:1:16 D
 . I V#2=1 S BS=$ZBITSET(BS,I,1)
 . S V=V\2
 Q BS
 ;
BS2I(BS) ; Convert 16-bit YottaDB $ZBITSTR back to integer
 N I,R,P S R=0,P=1
 F I=1:1:16 D
 . I $ZBITGET(BS,I) S R=R+P
 . S P=P*2
 Q R
 ;
BITAND(A,B) ; Bitwise AND via YottaDB $ZBITAND
 Q $$BS2I($ZBITAND($$I2BS(A),$$I2BS(B)))
 ;
BITOR(A,B) ; Bitwise OR via YottaDB $ZBITOR
 Q $$BS2I($ZBITOR($$I2BS(A),$$I2BS(B)))
 ;
BITXOR(A,B) ; Bitwise XOR via YottaDB $ZBITXOR
 Q $$BS2I($ZBITXOR($$I2BS(A),$$I2BS(B)))
 ;
BITCNT(A) ; Population count (1-bits) via YottaDB $ZBITCOUNT
 Q $ZBITCOUNT($$I2BS(A))
