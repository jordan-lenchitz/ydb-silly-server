ZTEST ; 🧪 NATIVE MUMPS ASSERTION & SELF-TEST SUITE
 W $$ALL(),!
 Q
 ;
ASSERT(NAME,GOT,WANT) ; Record test assertion in ^YDBTEST
 N ID S ID=$I(^YDBTEST("TOTAL"))
 I GOT=WANT D  Q:$Q 1  Q
 . S ^YDBTEST("PASS")=$G(^YDBTEST("PASS"))+1
 . S ^YDBTEST("CASES",ID)="PASS|"_NAME
 S ^YDBTEST("FAIL")=$G(^YDBTEST("FAIL"))+1
 S ^YDBTEST("CASES",ID)="FAIL|"_NAME_"|got="_GOT_"|want="_WANT
 Q:$Q 0
 Q
 ;
ALL() ; Run full native MUMPS test suite across all routines
 N N,V,A,OUT,LINES,JOBJ
 K ^YDBTEST
 S ^YDBTEST("PASS")=0,^YDBTEST("FAIL")=0,^YDBTEST("TOTAL")=0
 S ^YDBTEST("TS")=$ZUT
 ;
 ; 1. XEXE & XUTIL
 S N=$$XEXE^XEXE("W ""HELLO WORLD""")
 D ASSERT("XEXE.SingleLineNoBang",N,1)
 D ASSERT("XEXE.SingleLineVal",$G(^XOUT($J,1)),"HELLO WORLD")
 S N=$$XEXE^XEXE("W ""L1"",!,""L2"",!")
 D ASSERT("XEXE.MultiLine",N,2)
 D ASSERT("XEXE.MultiLine2",$G(^XOUT($J,2)),"L2")
 S N=$$XEXE^XEXE("S X=1/0")
 D ASSERT("XEXE.ErrorTrapCaptured",N>0,1)
 D ASSERT("XEXE.ErrorTrapPrefix",$E($G(^XOUT($J,1)),1,13),"[MUMPS-ERROR]")
 D ASSERT("XEXE.Eval",$$EVAL^XEXE("21*2"),42)
 D ASSERT("XEXE.VersionNonEmpty",$L($$VERSION^XEXE())>0,1)
 D ASSERT("XEXE.IntrinsicsHasZV",$$INTRINSICS^XEXE()["zv=",1)
 D ASSERT("XUTIL.getjob",$$getjob^XUTIL(),$J)
 D ASSERT("XUTIL.zyrel",$$zyrel^XUTIL()["YottaDB",1)
 D ASSERT("XUTIL.zblen",$$zblen^XUTIL("🚀"),4)
 D ASSERT("XUTIL.zut",$$zut^XUTIL()>0,1)
 ;
 ; 2. DATA (B*-Tree Persistence)
 S V=$$KILL^DATA("YTEST","")
 D ASSERT("DATA.Set",$$SET^DATA("YTEST","a|b|c","val123"),1)
 D ASSERT("DATA.Get",$$GET^DATA("YTEST","a|b|c"),"val123")
 D ASSERT("DATA.Exists",$$EXISTS^DATA("YTEST","a|b|c")>0,1)
 D ASSERT("DATA.Incr1",$$INCR^DATA("YTEST","counter",5),5)
 D ASSERT("DATA.Incr2",$$INCR^DATA("YTEST","counter",3),8)
 D ASSERT("DATA.Count",$$COUNT^DATA("YTEST"),2)
 D ASSERT("DATA.DumpHasVal",$$DUMP^DATA("YTEST",10)["val123",1)
 D ASSERT("DATA.Kill",$$KILL^DATA("YTEST",""),1)
 D ASSERT("DATA.ExistsAfterKill",$$EXISTS^DATA("YTEST","a|b|c"),0)
 ;
 ; 3. VMMGR (^YDBCLOUD Orchestrator)
 S V=$$CLEAN^VMMGR("vm-mumps-1")
 D ASSERT("VMMGR.Provision",$$PROVISION^VMMGR("vm-mumps-1","7"),"SUCCESS")
 D ASSERT("VMMGR.StatusProv",$$GETVM^VMMGR("vm-mumps-1","status"),"PROVISIONING")
 D ASSERT("VMMGR.CpuLimit",$$GETMETA^VMMGR("vm-mumps-1","cpuLimit"),"7.0")
 D ASSERT("VMMGR.Complete",$$COMPLETE^VMMGR("vm-mumps-1"),"1")
 D ASSERT("VMMGR.StatusRun",$$GETVM^VMMGR("vm-mumps-1","status"),"RUNNING")
 D ASSERT("VMMGR.Heartbeat",$$HEARTBEAT^VMMGR("vm-mumps-1")>0,1)
 D ASSERT("VMMGR.ListContains",$$LIST^VMMGR()["vm-mumps-1",1)
 D ASSERT("VMMGR.CountPos",$$COUNT^VMMGR()>0,1)
 D ASSERT("VMMGR.Stop",$$STOP^VMMGR("vm-mumps-1"),"1")
 D ASSERT("VMMGR.StatusStop",$$GETVM^VMMGR("vm-mumps-1","status"),"STOPPED")
 D ASSERT("VMMGR.Clean",$$CLEAN^VMMGR("vm-mumps-1"),"1")
 ;
 ; 4. LOG (^YDBLOGS Engine)
 S V=$$CLEAR^LOG("TESTLVL")
 D ASSERT("LOG.Write1",$$WRITE^LOG("TESTLVL","first msg"),1)
 D ASSERT("LOG.Write2",$$WRITE^LOG("TESTLVL","second msg"),2)
 D ASSERT("LOG.ReadHasMsg",$$READ^LOG("TESTLVL",2)["second msg",1)
 D ASSERT("LOG.StatsHasLvl",$$STATS^LOG()["TESTLVL:2;",1)
 D ASSERT("LOG.Clear",$$CLEAR^LOG("TESTLVL"),"1")
 ;
 ; 5. SYS ($Z* System Introspection)
 D ASSERT("SYS.Version",$L($$VERSION^SYS())>0,1)
 D ASSERT("SYS.PID",$$PID^SYS(),$J)
 D ASSERT("SYS.DateTimeLen",$L($$DATETIME^SYS()),19)
 D ASSERT("SYS.Health",$$HEALTH^SYS(),"OK")
 D ASSERT("SYS.MetricsHasPID",$$METRICS^SYS()[("pid="_$J),1)
 ;
 ; 6. MATH & $ZBIT* Bitwise Operations
 D ASSERT("MATH.Fact0",$$FACT^MATH(0),1)
 D ASSERT("MATH.Fact6",$$FACT^MATH(6),720)
 D ASSERT("MATH.IsPrime2",$$ISPRIME^MATH(2),1)
 D ASSERT("MATH.IsPrime97",$$ISPRIME^MATH(97),1)
 D ASSERT("MATH.IsPrime100",$$ISPRIME^MATH(100),0)
 D ASSERT("MATH.GCD",$$GCD^MATH(48,18),6)
 D ASSERT("MATH.LCM",$$LCM^MATH(12,18),36)
 D ASSERT("MATH.Abs",$$ABS^MATH(-42),42)
 D ASSERT("MATH.Pow",$$POW^MATH(2,10),1024)
 D ASSERT("MATH.Round",$$ROUND^MATH(3.14159,2),3.14)
 D ASSERT("MATH.Fib10",$$FIB^MATH(10),55)
 D ASSERT("MATH.DH",$$DH^MATH(255),"FF")
 D ASSERT("MATH.HD",$$HD^MATH("ff"),255)
 D ASSERT("MATH.BitAnd",$$BITAND^MATH(12,10),8)
 D ASSERT("MATH.BitOr",$$BITOR^MATH(12,10),14)
 D ASSERT("MATH.BitXor",$$BITXOR^MATH(12,10),6)
 D ASSERT("MATH.BitCnt",$$BITCNT^MATH(255),8)
 ;
 ; 7. STR (UTF-8 & Byte String Operations)
 D ASSERT("STR.Rev",$$REV^STR("gopher"),"rehpog")
 D ASSERT("STR.Up",$$UP^STR("ydb-silly"),"YDB-SILLY")
 D ASSERT("STR.Low",$$LOW^STR("MUMPS"),"mumps")
 D ASSERT("STR.Trim",$$TRIM^STR("  hello  "),"hello")
 D ASSERT("STR.Replace",$$REPLACE^STR("a-b-c","-","_"),"a_b_c")
 D ASSERT("STR.LenUTF8",$$LEN^STR("🚀a"),2)
 D ASSERT("STR.BLenUTF8",$$BLEN^STR("🚀a"),5)
 D ASSERT("STR.Sub",$$SUB^STR("hello",2,4),"ell")
 D ASSERT("STR.Split",$$SPLIT^STR("x,y,z",",","A"),3)
 D ASSERT("STR.Join",$$JOIN^STR("A","-"),"x-y-z")
 ;
 ; 8. WEB (URL & Header & JSON)
 D ASSERT("WEB.URLEncDec",$$URLDEC^WEB($$URLENC^WEB("hello 🚀!")),"hello 🚀!")
 D ASSERT("WEB.HdrGet",$$HDRGET^WEB("Host: local"_$C(10)_"X-VM-ID: vm-99","x-vm-id"),"vm-99")
 K JOBJ S JOBJ("k")="v"
 D ASSERT("WEB.JSON",$$JSON^WEB("JOBJ"),"{""k"":""v""}")
 ;
 ; 9. ZDB, ZFORGE, ZREQ
 D ASSERT("ZDB.Seed",$$SEED^ZDB(),3)
 D ASSERT("ZDB.Add",$$ADD^ZDB("NYAN","POP-TART"),4)
 D ASSERT("ZDB.Count",$$COUNT^ZDB(),4)
 D ASSERT("ZDB.CatsJSON",$$C^ZDB()["MR. WHISKERS",1)
 S LINES=$$FORGE^ZFORGE()
 D ASSERT("ZFORGE.LinesPos",LINES>100,1)
 D ASSERT("ZFORGE.DecodeA",$$DECODE^ZFORGE("A")["ASCII_HEX_41",1)
 D ASSERT("ZFORGE.DecodeEmoji",$$DECODE^ZFORGE($C(128512))["EMOJI_HEX_1F600",1)
 D ASSERT("ZREQ.Wisdom",$$ROUTE^ZREQ("/api/wisdom")[("MUMPS IS ETERNAL"),1)
 D ASSERT("ZREQ.Dogs",$$ROUTE^ZREQ("/api/dogs")[("MAXIMUM BORK"),1)
 D ASSERT("ZREQ.Cats",$$ROUTE^ZREQ("/api/cats")[("GARFIELD"),1)
 ;
 Q "PASS:"_^YDBTEST("PASS")_"|FAIL:"_^YDBTEST("FAIL")_"|TOTAL:"_^YDBTEST("TOTAL")
