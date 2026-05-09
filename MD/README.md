# 🚀 YottaDB Native Proxy V2.0.0

## "Way Better, Fast, and Good"

This is a high-performance Node.js proxy for YottaDB, re-engineered for unironic speed and production-grade reliability.

### ⚡ Key Improvements
- **Native N-API Binding:** Switched from `child_process.spawn` to `mg-dbx-napi` for all database interactions. Execution is now handled in-process, bypassing shell overhead.
- **Optimized Executioner:** Uses a custom MUMPS routine (`XEXE.m`) to capture `stdout` via temporary stream buffers, allowing for complex MUMPS code execution with full output capture.
- **Persistence:** VM states and system metadata are now persisted directly in YottaDB globals (`^YDBCLOUD`), ensuring state survives proxy restarts.
- **Extended API:** 
    - `POST /api/execute`: Native high-speed MUMPS execution.
    - `GET /api/global/:name`: Direct global inspection.
    - `POST /api/global/:name`: Direct global manipulation.
    - `GET /api/vm/status`: Real-time engine and VM health.
- **Isolation Ready:** Architected to support multi-region/multi-tenant global directory isolation.

### 🛠️ Architecture
The proxy acts as a bridge between modern RESTful microservices and the raw power of YottaDB. By using the `mg-dbx-napi` N-API wrapper, we achieve near-memory speeds for database operations while maintaining the flexibility of Express.js.

### 🏃 Running
```bash
npm install
npm start
```

### ✅ Bonafides Verified
- [x] Zero process spawning on hot paths.
- [x] Native MUMPS routine integration.
- [x] Global-backed state management.
- [x] Optimized I/O streams.
