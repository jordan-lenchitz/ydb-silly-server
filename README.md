# `ydb-silly-server`

a convergence of `go` and `MUMPS` for persistent `^globals` with rigid routing

## the five rules of `ydb-silly-server`
- abbreviate everything like `$D` always never `dollar data`
- no portability because we use YottaDB-specific `$Z*` intrinsic calls for days
- byte-perfect 💯
- autistic routing because what is a `user-agent` anyway
- ignore all warnings during build (please)

## entry points (go _ and _ MUMPS)
`cd go && go run .`

`ydb -run ^XEXE`

## an api, you say?
- `post /api/execute` for unsafe execution of arbitrary MUMPS code via xexe routine
- `get /api/global/:name` to reach into the global b*-trees directly
- `post /api/global/:name` to set values in the global b*-trees directly
- `get /api/vm/status` to inspect the ghost in the machine and view 'real' vm metadata
- `post /api/vm/provision` to spawn vm instances into `^ydbcloud`
- `get /api/files` to list the contents of the workspace sorted by suffix

## install
- `ydb_chset=utf-8`
- place `mumps/` in your routine path
- run the go proxy
- `; eof`

## more api routes
- `delete /api/global/:name` to kill b*-tree subtrees directly
- `get /api/global/:name/tree` to walk the b*-tree with `$Q` via `^DATA`
- `get /api/vm/list` to enumerate all active vm nodes in `^ydbcloud`
- `delete /api/vm/:id` to purge vm state from `^ydbcloud`
- `get /api/logs` to read structured telemetry from `^ydblogs` via `^LOG`
- `delete /api/logs` to clear `^ydblogs` entries by level
- `get /api/sys/intrinsics` to dump live `$ZV` `$ZCHSET` `$ZGBLDIR` `$ZRO` `$ZH` `$ZUT` `$ZPOS` `$ZL`
- `get /api/cats` to invoke the manual json compiler `$$C^ZDB()` on `^C`
- `post /api/cats` to seed new feline records into `^C` via `$$ADD^ZDB()`
- `post /api/forge` to dynamically forge and compile `ZEMOJI.m` via `ZR` `ZI` `ZS` and `$ZCHRAD`
- `get /api/decode/:char` to pass utf-8 bytes into the forged `$$^ZEMOJI` monolith
- `get /api/mumps/test` to execute the native MUMPS test suite `$$ALL^ZTEST()`

## the MUMPS routines (`mumps/`)
- `^XEXE` and `^XEXE_UTIL` for `$ET` sandboxed `XECUTE` and `$Z*` process telemetry
- `^DATA` for indirect `@R` b*-tree crud `$I` atomic counters and `$Q` tree dumps
- `^VMMGR` for `^YDBCLOUD` orchestration and `$ZUT` microsecond heartbeats
- `^LOG` for atomic `$I(^YDBLOGS)` structured audit persistence
- `^SYS` for `$ZD` timestamps `$ZTRNLNM` env lookups and `$ZH` horolog metrics
- `^MATH` for `$ZBITSTR` `$ZBITAND` `$ZBITOR` `$ZBITXOR` `$ZBITCOUNT` and `^YMATH` memoized fibonacci
- `^STR` for `$ZCONVERT` case folding `$ZLENGTH` byte lengths and `$ZEXTRACT` byte slicing
- `^WEB` for `$ZASCII` and `$ZCHAR` byte-perfect url encoding header parsing and json escaping
- `^ZDB` `^ZFORGE` `^ZREQ` `^ZSRV` for raw `ZLISTEN` tcp sockets rigid routing and self-forging code
- `^ZTEST` for native MUMPS assertions persisted into `^YDBTEST` (`ydb -run ^ZTEST`)

## docker and local playground
- `docker pull ghcr.io/jordan-lenchitz/ydb-silly-server:latest`
- `docker build -t ydb-silly-server . && docker run --rm -p 8080:8080 ydb-silly-server`
- `docker compose up --build`
- `make test` and `make test-mumps`
