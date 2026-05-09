#!/bin/bash
# Source YottaDB environment
. /opt/yottadb/current/ydb_env_set

# Manually export ydb_dist if it's missing (it usually is after sourcing ydb_env_set in some shells)
export ydb_dist=${ydb_dist:-/opt/yottadb/current}

# Add local directory to routines
export ydb_routines=". /app $ydb_routines"

echo "--- YDB ENV ---"
env | grep ydb
echo "---------------"

# Initialize database directory if missing
# Note: On Cloud Run, /data might be local and ephemeral if not mounted
GBL_DIR=$(dirname "$ydb_gbldir")
if [ ! -d "$GBL_DIR" ]; then
    echo "Initializing YottaDB database directory $GBL_DIR..."
    mkdir -p "$GBL_DIR"
    mupip create
else
    echo "Database exists. Running recovery..."
    mupip rundown -region "*"
    mupip journal -recover -backward "*"
    mupip rundown -region "*"
fi

# Start Node.js application
echo "Starting Node.js application..."
npm start
