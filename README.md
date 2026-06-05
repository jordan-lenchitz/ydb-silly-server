# yottadb silly server

the convergence of go and mumps for persistent globals and rigid routing

## the tenets
1 abbreviate everything
2 no portability
3 byte perfect
4 autistic routing
5 ignore all warnings

## entry points
cd go && go run .
ydb -run ^XEXE

## the api
post /api/execute
unsafe execution of arbitrary mumps code via xexe routine

get /api/global/:name
reach into the global b trees directly

post /api/global/:name
set values in the global b trees directly

get /api/vm/status
inspect the ghost in the machine and view fake vm metadata

post /api/vm/provision
spawn fake instances into the ydbcloud global

get /api/files
list the contents of the workspace sorted by suffix

## install
1 ydb_chset=utf-8
2 place mumps/ in your routine path
3 run the go proxy

; eof
