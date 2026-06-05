#!/bin/bash
# source yottadb environment
. /opt/yottadb/current/ydb_env_set

# manually export ydb_dist if it is missing
export ydb_dist=${ydb_dist:-/opt/yottadb/current}

# add local directory to routines
export ydb_routines=". ./MUMPS $ydb_routines"

echo "--- ydb env ---"
env | grep ydb
echo "---------------"

# initialize database directory if missing
GBL_DIR=$(dirname "$ydb_gbldir")
if [ ! -d "$GBL_DIR" ]; then
    echo "initializing YottaDB database directory $GBL_DIR..."
    mkdir -p "$GBL_DIR"
    mupip create
else
    echo "database exists; running recovery..."
    mupip rundown -region "*"
    mupip journal -recover -backward "*" || true
    mupip rundown -region "*"
fi

# start the application
echo "starting go application..."
./ydb-server
