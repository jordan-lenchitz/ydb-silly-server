# 🚀 yottadb silly server (YSS)

100% silly high-performance Node.js proxy for YottaDB

### ⚡ key features
- uses `mg-dbx-napi` for in-process database interactions
- uses `MUMPS/XEXE.m` to capture output via stream buffers
- persists state in YottaDB globals (`^YDBCLOUD`)
- offers an API
    - `POST /api/execute` = native high-speed MUMPS execution
    - `GET /api/global/:name` = direct global inspection
    - `POST /api/global/:name` = direct global manipulation
    - `GET /api/vm/status` = real-time engine and VM health

### 🛠️ architecture
by using the `mg-dbx-napi` N-API wrapper we achieve near-memory speeds for database operations while maintaining the flexibility of express.js

### 📂 structure
- `JS/` = node.js application logic
- `MUMPS/` = native YottaDB routines
- `ydb-data/` = persistent database storage (untracked, managed by Docker volumes)
- `logs/` = application and installation logs (untracked, generated at runtime)

### 🏃 running
```bash
npm install
npm start
```

### ✅ bonafides 
- [x] 0 process spawning on hot paths
- [x] native MUMPS routine integration
- [x] global-backed state management
