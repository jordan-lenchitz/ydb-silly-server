import express from 'express';
import bodyParser from 'body-parser';
import { server, mglobal } from 'mg-dbx-napi';
import { spawn } from 'child_process';
import fs from 'fs';
import path from 'path';

/**
 * 🚀 YottaDB UNIRONIC BACKEND V2.1.0
 * Blazing fast, Bloat-free, Folders for all!
 */

const app = express();
app.use(bodyParser.json());

const BASE_DIR = '/app';
const FOLDERS = ['MUMPS', 'JS', 'MD', 'logs'];

// CORS - The "Good" Way
app.use((req, res, next) => {
    res.header("Access-Control-Allow-Origin", "*");
    res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, X-VM-ID");
    res.header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS");
    if (req.method === 'OPTIONS') return res.sendStatus(200);
    next();
});

class YottaDBEngine {
    constructor() {
        this.db = new server();
        this.config = {
            type: "YottaDB",
            path: "/opt/yottadb/current",
            env_vars: {
                ydb_gbldir: process.env.ydb_gbldir || "/data/r2.06_x86_64/g/yottadb.gld",
                ydb_routines: process.env.ydb_routines || ". /app/MUMPS /app /opt/yottadb/current"
            }
        };
        this.connected = false;
        this.jobId = null;
    }

    init() {
        console.log("--- ⚡ BOOTING YOTTADB NATIVE ENGINE ---");
        try {
            this.db.open(this.config);
            this.connected = true;
            this.jobId = this.db.function("getjob^XEXE_UTIL");
            console.log(`[ENGINE] Connected. Version: ${this.db.dbversion()}`);
            console.log(`[ENGINE] Native Process ID: ${this.jobId}`);
        } catch (err) {
            console.error("[ENGINE] Connection Failed. Falling back to CLI mode.", err);
        }
    }

    async execute(mCode) {
        if (!this.connected) return this.executeCLI(mCode);
        try {
            const lineCount = this.db.function("XEXE^XEXE", mCode);
            let output = "";
            for (let i = 1; i <= lineCount; i++) {
                const line = this.db.get("XOUT", [this.jobId, i]);
                output += (line || "") + "\n";
            }
            return { output: output.trim(), status: 'OK', method: 'NATIVE' };
        } catch (err) {
            console.error("[ENGINE] Native Execution Error:", err);
            return { error: err.message, status: 'FAIL' };
        }
    }

    executeCLI(mCode) {
        return new Promise((resolve) => {
            const ydb = spawn('/opt/yottadb/current/yottadb', ['-run', '%XCMD', mCode], {
                env: process.env
            });
            let stdout = '', stderr = '';
            ydb.stdout.on('data', d => stdout += d);
            ydb.stderr.on('data', d => stderr += d);
            ydb.on('close', code => resolve({
                output: (stdout || stderr).trim(),
                status: code === 0 || stdout ? 'OK' : 'FAIL',
                method: 'CLI'
            }));
        });
    }
}

const engine = new YottaDBEngine();
engine.init();

// Persistence Helpers
const getVMState = (id) => {
    if (!id) return null;
    try {
        const state = {
            status: engine.db.function("GETVM^VMMGR", id, "status"),
            instanceId: engine.db.function("GETVM^VMMGR", id, "instanceId"),
            provisionedAt: engine.db.function("GETVM^VMMGR", id, "provisionedAt"),
            provisioned: engine.db.function("GETVM^VMMGR", id, "provisioned") === "1",
            metadata: {
                internalIp: engine.db.function("GETMETA^VMMGR", id, "internalIp"),
                port: parseInt(engine.db.function("GETMETA^VMMGR", id, "port")),
                region: engine.db.function("GETMETA^VMMGR", id, "region"),
                cpuLimit: engine.db.function("GETMETA^VMMGR", id, "cpuLimit")
            }
        };
        return state.status ? state : null;
    } catch { return null; }
};

// API Endpoints

app.get('/api/vm/status', (req, res) => {
    const vmId = req.headers['x-vm-id'];
    const state = getVMState(vmId) || { status: 'OFFLINE', provisioned: false };
    res.send({
        ...state,
        system: {
            engine: engine.db.dbversion(),
            nativeJob: engine.jobId,
            uptime: engine.db.function("UPTIME^SYS")
        }
    });
});

app.post('/api/vm/provision', (req, res) => {
    const newId = `ydb-${Math.random().toString(36).substring(2, 10)}`;
    const cores = [3, 5, 7][Math.floor(Math.random() * 3)];

    // Call native MUMPS provisioner
    engine.db.function("PROVISION^VMMGR", newId, cores.toString());
    engine.db.function("INFO^LOG", `Initiated provisioning for ${newId}`);

    const state = getVMState(newId);
    
    setTimeout(() => {
        engine.db.function("COMPLETE^VMMGR", newId);
        engine.db.function("INFO^LOG", `Provisioning complete for ${newId}`);
    }, 2000);

    res.send({ message: 'Provisioning initiated', instanceId: newId, ...state });
});

app.post(['/execute', '/api/execute'], async (req, res) => {
    const { mCode } = req.body;
    if (!mCode) return res.status(400).send({ error: 'No code provided' });
    const result = await engine.execute(mCode);
    res.send(result);
});

/**
 * 📂 FILE EXPLORER API
 * "Sort by file type suffix ONLY"
 */
app.get('/api/files', (req, res) => {
    let allFiles = [];

    FOLDERS.forEach(folder => {
        const dirPath = path.join(BASE_DIR, folder);
        if (fs.existsSync(dirPath)) {
            const files = fs.readdirSync(dirPath).map(f => ({
                name: f,
                folder: folder,
                path: path.join(folder, f),
                suffix: path.extname(f).toLowerCase()
            }));
            allFiles = allFiles.concat(files);
        }
    });

    // Add root files
    fs.readdirSync(BASE_DIR).forEach(f => {
        const fullPath = path.join(BASE_DIR, f);
        if (fs.statSync(fullPath).isFile()) {
            allFiles.push({
                name: f,
                folder: '/',
                path: f,
                suffix: path.extname(f).toLowerCase()
            });
        }
    });

    // Sort by suffix ONLY
    allFiles.sort((a, b) => {
        if (a.suffix < b.suffix) return -1;
        if (a.suffix > b.suffix) return 1;
        return a.name.localeCompare(b.name);
    });

    res.send({ files: allFiles });
});

app.get('/api/global/:name', (req, res) => {
    const { name } = req.params;
    const subs = req.query.subs ? req.query.subs.split(',') : [];
    try {
        const val = engine.db.get(name, subs);
        res.send({ global: name, subscripts: subs, value: val });
    } catch (err) {
        res.status(500).send({ error: err.message });
    }
});

app.post('/api/global/:name', (req, res) => {
    const { name } = req.params;
    const { subscripts, value } = req.body;
    try {
        engine.db.set(name, subscripts || [], value);
        res.send({ status: 'OK' });
    } catch (err) {
        res.status(500).send({ error: err.message });
    }
});

const PORT = process.env.PORT || 8080;
app.listen(PORT, () => {
    console.log(`
    =========================================
    🚀 YOTTADB NATIVE PROXY V2.1.0
    PORT: ${PORT}
    MODE: NATIVE (N-API)
    BONAFIDES: FOLDERIZED ✅
    =========================================
    `);
});
