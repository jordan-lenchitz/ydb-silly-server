ZDB ; 💾 THE DATA CRYPT
SEED ; 🪴 PLANT THE SILLY SEEDS
 K ^C
 S ^C(1,"NAME")="MR. WHISKERS",^C(1,"MOOD")="HISSING"
 S ^C(2,"NAME")="GARFIELD",^C(2,"MOOD")="LASAGNA"
 S ^C(3,"NAME")="BINGUS",^C(3,"MOOD")="BALD"
 Q
C() ; 🐈 THE AUTISTICALLY MANUAL JSON COMPILER
 N I,R,K,V,F1,F2
 S R="{""cats"":[",I=0,F1=1
L1 S I=$O(^C(I)) G:I="" D
 S R=R_$S(F1:"",1:",")_"{"
 S K="",F1=0,F2=1
L2 S K=$O(^C(I,K)) G:K="" L3
 S V=^C(I,K)
 S R=R_$S(F2:"",1:",")_""""_K_""":"""_V_""""
 S F2=0 G L2
L3 S R=R_"}" G L1
D S R=R_"]}"
 Q R
