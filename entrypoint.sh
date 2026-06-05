#!/bin/bash
# Source YottaDB environment
. /opt/yottadb/current/ydb_env_set

# Manually export ydb_dist if it's missing
export ydb_dist=${ydb_dist:-/opt/yottadb/current}

# Add local directory to routines
export ydb_routines=". /app/MUMPS /app $ydb_routines"

echo "--- YDB ENV ---"
env | grep ydb
echo "---------------"

# Initialize database directory if missing
GBL_DIR=$(dirname "$ydb_gbldir")
if [ ! -d "$GBL_DIR" ]; then
    echo "Initializing YottaDB database directory $GBL_DIR..."
    mkdir -p "$GBL_DIR"
    mupip create
else
    echo "Database exists. Running recovery..."
    mupip rundown -region "*"
    mupip journal -recover -backward "*" || true
    mupip rundown -region "*"
fi

# Start Go application
echo "Starting Go application..."
./ydb-server
