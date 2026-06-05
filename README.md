# yottadb-silly-server

a silly convergence of golang and MUMPS!

persistent nosql storage with rigid routing and absolutely zero abstraction :)

## the tenets
1. **abbreviate everything** 
2. **no portability** 
3. **byte-perfect** 
4. **the forge of routines that build routines**
5. **autistic routing** 

## entry points
```bash
# the go layer
cd go && go run .

# the MUMPS layer
ydb -run ^ZSRV
```

## an api, you say?
- `POST /api/execute { "mCode": "..." }` is for unsafe execution of arbitrary MUMPS code
- `GET /api/global/:name` hikes into the `B*-tree` forest so you do not have to

## howto (install)
1. `ydb_chset=utf-8`
2. `place MUMPS/` in your routine path
3. `D ^ZFORGE`
4. ignore all warnings as they are the best part of the process!

; EOF
