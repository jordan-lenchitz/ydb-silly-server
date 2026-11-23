XUTIL ; ⚙️ YOTTADB $Z* PROCESS & BYTE INTRINSIC UTILITIES
 ;
getjob() Q $J
 ;
zver() Q $ZV
 ;
zyrel() Q $ZYRELEASE
 ;
zro() Q $ZRO
 ;
zgbldir() Q $ZGBLDIR
 ;
zchset() Q $ZCHSET
 ;
zut() Q $ZUT
 ;
zhorolog() Q $ZH
 ;
zlevel() Q $ZL
 ;
zpos() Q $ZPOS
 ;
ztrn(nm) Q $ZTRNLNM($G(nm))
 ;
zparse(f,fld) Q $ZPARSE($G(f),$G(fld))
 ;
zblen(s) Q $ZLENGTH($G(s))
 ;
zbsub(s,a,b) Q $ZEXTRACT($G(s),+$G(a),+$G(b))
 ;
zbpiece(s,d,a) Q $ZPIECE($G(s),$G(d),+$G(a))
