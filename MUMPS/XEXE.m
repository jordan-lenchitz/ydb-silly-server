XEXE(CODE) ; 🚀 THE UNHOLY EXECUTIONER
 NEW %FN,%LINE,%OUT,%I
 SET %OUT=$J
 KILL ^XOUT(%OUT)
 SET %FN="/tmp/ydb_out_"_%OUT
 OPEN %FN:(newversion:truncate:stream):0
 USE %FN
 XECUTE CODE
 CLOSE %FN
 OPEN %FN:(readonly:stream):0
 SET %I=0
 FOR  USE %FN READ %LINE QUIT:$ZEOF  SET %I=%I+1,^XOUT(%OUT,%I)=%LINE
 CLOSE %FN
 ZSYSTEM "rm "_%FN
 QUIT %I ; Return number of lines
