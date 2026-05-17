# 🚀 YottaDB Native Proxy V2.1.0

## "Folderized for Your Pleasure"

This is a high-performance Node.js proxy for YottaDB, re-engineered for unironic speed and folderized for production-grade organization.

### ⚡ Key Improvements
- **Folderized Architecture:** Clean root directory with logic separated into `JS/`, `MUMPS/`, and `logs/`.
- **File Explorer API:** New `GET /api/files` endpoint to navigate the project structure.
- **Native N-API Binding:** Using `mg-dbx-napi` for in-process database interactions.
- **Optimized Executioner:** Uses `MUMPS/XEXE.m` to capture output via stream buffers.
- **Persistence:** State persisted in YottaDB globals (`^YDBCLOUD`).
- **Extended API:** 
    - `POST /api/execute`: Native high-speed MUMPS execution.
    - `GET /api/global/:name`: Direct global inspection.
    - `POST /api/global/:name`: Direct global manipulation.
    - `GET /api/vm/status`: Real-time engine and VM health.
    - `GET /api/files`: Metadata-driven file exploration.

### 🛠️ Architecture
The proxy acts as a bridge between modern RESTful microservices and the raw power of YottaDB. By using the `mg-dbx-napi` N-API wrapper, we achieve near-memory speeds for database operations while maintaining the flexibility of Express.js.

### 📂 Structure
- `JS/`: Node.js application logic.
- `MUMPS/`: Native YottaDB routines.
- `ydb-data/`: Persistent database storage.
- `logs/`: Application and installation logs.

### 🏃 Running
```bash
npm install
npm start
```

### ✅ Bonafides Verified
- [x] Folderized for maximum cleanliness.
- [x] Zero process spawning on hot paths.
- [x] Native MUMPS routine integration.
- [x] Global-backed state management.
