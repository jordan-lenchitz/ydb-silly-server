ZFORGE ; 🔨 THE UNHOLY ROUTINE FORGE
 N %L,%I,%C,%H,%S,%U
 W #,"[FORGING 100,000 LINES OF STUPIDITY...]",!
 S %L(1)="ZEMOJI(%E) ; 🗿 THE MEGA-MATRIX OF ALL HUMAN KNOWLEDGE"
 S %L(2)=" S %B="""",""%N="""",""%M="""""
 S %L(3)=" ; ⚠️ WARNING: THIS ROUTINE WAS FORGED BY A ROBOT"
 S %C=4
 ; 🌀 LOOP THROUGH EMOJI HEX BLOCKS (U+1F600 to U+1F64F)
 F %I=128512:1:128591 D
 . S %H=$$DH(%I)
 . S %U=$ZCHRAD(%I) ; YOTTADB NATIVE UTF8 DECODE
 . S %L(%C)=" I %E="""_%U_""" D  G Q",%C=%C+1
 . S %L(%C)=" . S %N=""EMOJI_HEX_"_%H_"""",%C=%C+1
 . S %L(%C)=" . S %M=""REEEEEEEEEE""",%C=%C+1
 . S %L(%C)=" . S %B=""{""""char"""":""""""_%E_"""""",""""hex"""":""""""_%N_"""""",""""mood"""":""""""_%M_""""""}""",%C=%C+1
 ; 🌀 LOOP THROUGH BASIC LATIN/UTF8 JUST TO MAKE IT FAT
 F %I=33:1:126 D
 . S %U=$C(%I),%H=$$DH(%I)
 . S %L(%C)=" I %E="""_%U_""" D  G Q",%C=%C+1
 . S %L(%C)=" . S %N=""ASCII_HEX_"_%H_"""",%C=%C+1
 . S %L(%C)=" . S %M=""BORING_TEXT""",%C=%C+1
 . S %L(%C)=" . S %B=""{""""char"""":""""""_%E_"""""",""""hex"""":""""""_%N_"""""",""""mood"""":""""""_%M_""""""}""",%C=%C+1
 ; 🏁 THE END OF THE FILE
 S %L(%C)="Q ; 🛑 CATCH-ALL EXIT"
 S %L(%C+1)=" . S %B=""{""""error"""":""""UNKNOWN ENTITY""""}"""
 S %L(%C+2)=" Q %B"
 ; 💾 INJECT DIRECTLY INTO YOTTADB COMPILER
 ZR  F %S=1:1:%C+2 ZI %L(%S)
 ZS ZEMOJI
 W "[FORGE COMPLETE. ZEMOJI.m HAS BEEN BORN]",!
 Q
DH(%D) ; 🧮 DECIMAL TO HEX CONVERTER (AUTISTICALLY MANUAL)
 N %X,%R,%Y S %X="0123456789ABCDEF",%R=""
 I %D=0 Q "0"
 F  Q:%D=0  S %Y=%D#16,%R=$E(%X,%Y+1)_%R,%D=%D\16
 Q %R
