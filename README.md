# yss (yottadb-silly-server)

the unholy convergence of go and m. persistent globals. rigid routing. zero abstraction.

## the tenets
1. **abbreviate everything.** $D is not $DATA. $O is not $ORDER. if a routine name has more than 8 characters, it is bloat.
2. **no portability.** we exploit zlisten and job. we live in the distribution-specific dirt.
3. **byte-perfect.** exact boundaries. $C(13,10). no middleware to save you.
4. **the forge.** routines that build routines. code that spawns code.
5. **autistic routing.** user-agents are irrelevant. only the byte stream remains.

## entry points
```bash
# the gopher layer
cd go && go run .

# the m layer
ydb -run ^ZSRV
```

## api.txt
- `GET /api/vm/status` -> return the ghost in the machine.
- `POST /api/execute` -> { "mCode": "..." } -> unsafe execution. do not use.
- `GET /api/global/:name` -> reach into the global trees.

## install
1. ydb_chset=utf-8
2. place MUMPS/ in your routine path
3. D ^ZFORGE
4. ignore the warnings. they are part of the process.

; EOF
