ZDB ; 💾 THE DATA CRYPT (^C CATS PERSISTENCE)
 ;
SEED() ; 🪴 PLANT THE SILLY SEEDS
 K ^C
 S ^C(1,"NAME")="MR. WHISKERS",^C(1,"MOOD")="HISSING"
 S ^C(2,"NAME")="GARFIELD",^C(2,"MOOD")="LASAGNA"
 S ^C(3,"NAME")="BINGUS",^C(3,"MOOD")="BALD"
 Q:$Q 3
 Q
 ;
ADD(NAME,MOOD) ; 🐾 ADD A NEW CAT TO THE CRYPT
 N ID
 I '$D(^C) D SEED()
 S ID=+$O(^C(""),-1)+1
 S ^C(ID,"NAME")=$S($G(NAME)'="":NAME,1:"MYSTERY CAT")
 S ^C(ID,"MOOD")=$S($G(MOOD)'="":MOOD,1:"CHAOTIC")
 Q:$Q ID
 Q
 ;
COUNT() ; 🧮 COUNT CATS IN ^C
 N I,C S I=0,C=0
 F  S I=$O(^C(I)) Q:I=""  S C=C+1
 Q C
 ;
C() ; 🐈 THE AUTISTICALLY MANUAL JSON COMPILER
 N I,R,K,V,F1,F2
 I '$D(^C) D SEED()
 S R="{""cats"":[",I=0,F1=1
L1 S I=$O(^C(I)) G:I="" D
 S R=R_$S(F1:"",1:",")_"{"
 S K="",F1=0,F2=1
L2 S K=$O(^C(I,K)) G:K="" L3
 S V=$$JSONESC^WEB($G(^C(I,K)))
 S R=R_$S(F2:"",1:",")_""""_$$JSONESC^WEB(K)_""":"""_V_""""
 S F2=0 G L2
L3 S R=R_"}" G L1
D S R=R_"]}"
 Q R
