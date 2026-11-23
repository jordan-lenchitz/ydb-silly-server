#!/usr/bin/env bash
set -e

YDB_DIST_DIR="${ydb_dist:-/opt/yottadb/current}"
if [ ! -f "$YDB_DIST_DIR/ydb_env_set" ]; then
    YDB_DIST_DIR=$(find /usr/local/lib/yottadb /opt/yottadb -name ydb_env_set -exec dirname {} \; 2>/dev/null | head -n 1)
fi

if [ -z "$YDB_DIST_DIR" ] || [ ! -f "$YDB_DIST_DIR/ydb_env_set" ]; then
    echo "[ERROR] Could not locate ydb_env_set"
    exit 1
fi

export ydb_chset="UTF-8"
export ydb_dir="${ydb_dir:-/data}"
mkdir -p "$ydb_dir"

# Source YottaDB environment (sets ydb_dist, ydb_gbldir, ydb_routines, etc.)
. "$YDB_DIST_DIR/ydb_env_set"

# Recover/rundown existing database files if container hostname changed or after unclean stop
for dat in "$ydb_dir"/*/g/*.dat; do
    if [ -f "$dat" ]; then
        "$ydb_dist/mupip" rundown -file "$dat" 2>/dev/null || true
    fi
done
"$ydb_dist/mupip" rundown -region "*" 2>/dev/null || true

APP_DIR="$(cd "$(dirname "$0")" && pwd)"
export ydb_routines="$APP_DIR/mumps $ydb_routines"

echo "--- ⚡ YOTTADB ENVIRONMENT ---"
env | grep -i ydb || true
echo "------------------------------"

echo "[BOOT] Running native MUMPS self-test suite (^ZTEST)..."
"$ydb_dist/yottadb" -run ^ZTEST

if [ "$1" = "test" ]; then
    cd "$APP_DIR/go" && go test -v ./...
    exit 0
fi

if [ -x "$APP_DIR/ydb-server" ]; then
    exec "$APP_DIR/ydb-server"
else
    cd "$APP_DIR/go" && exec go run .
fi
