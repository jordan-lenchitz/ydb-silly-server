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
- `post /api/execute` for unsafe execution of arbitrary mumps code via xexe routine
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
