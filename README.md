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

## the trials
- test vm status endpoint: verifies virtual machine status response structure
- test execute nil safety: ensures code execution handles missing database
- test concurrent database access: validates thread safety under heavy load
- test log invariant: checks logger never panics on input
- test factorial: confirms mathematical correctness for factorial results
- test is prime: validates prime number detection logic accuracy
- test gcd: verifies greatest common divisor calculation logic
- test lcm: ensures least common multiple calculation accuracy
- test fibonacci: confirms correct sequence generation for fibonacci
- test reverse: validates string reversal with utf8 characters
- test substring: checks rune based safe string slicing logic
- fuzz reverse: stress tests string reversal with random bytes
- fuzz substring: validates slicing logic against arbitrary inputs
- fuzz gcd: confirms mathematical properties for random integers
- fuzz url: ensures perfect round trip url encoding
- fuzz parse headers: stress tests header parser against slop
- benchmark factorial: measures raw performance of mathematical logic
- benchmark reverse: tracks memory allocations during string reversal
- benchmark reverse parallel: analyzes scaling across multiple cpu cores
- benchmark substring: measures throughput of safe string slicing
- benchmark gcd: evaluates efficiency of euclidean algorithm implementation

; EOF
