#!/usr/bin/env bash
# Local end-to-end test script for ydb-silly-server on Apple Silicon Macs (Colima / linux/arm64)
#
# Usage:
#   ./scripts/test-arm64-colima.sh              # Pull ghcr.io image & run full ARM64 test suite
#   ./scripts/test-arm64-colima.sh --build      # Build local Dockerfile for linux/arm64 & test
#   ./scripts/test-arm64-colima.sh --port 18080 # Use custom host port (default: 18080)
#   ./scripts/test-arm64-colima.sh --keep       # Keep container running after tests pass
#   ./scripts/test-arm64-colima.sh --arch amd64 # Override expected arch (default: arm64)

set -euo pipefail

IMAGE="${IMAGE:-ghcr.io/jordan-lenchitz/ydb-silly-server:latest}"
MODE="pull"
PORT="18080"
KEEP=0
EXPECT_ARCH="arm64"

while [ $# -gt 0 ]; do
    case "$1" in
        --build)
            MODE="build"
            IMAGE="ydb-silly-server:local-arm64"
            ;;
        --pull)
            MODE="pull"
            ;;
        --port)
            PORT="$2"
            shift
            ;;
        --image)
            IMAGE="$2"
            shift
            ;;
        --arch)
            EXPECT_ARCH="$2"
            shift
            ;;
        --keep)
            KEEP=1
            ;;
        -h|--help)
            sed -n '2,10p' "$0" | sed 's/^# \{0,1\}//'
            exit 0
            ;;
        *)
            echo "Unknown argument: $1" >&2
            exit 1
            ;;
    esac
    shift
done

if [ "$EXPECT_ARCH" = "arm64" ]; then
    EXPECT_UNAME="aarch64"
else
    EXPECT_UNAME="x86_64"
fi

CONTAINER_NAME="ydb-silly-test-$$"
VOLUME_NAME="ydb-silly-test-data-$$"
BASE_URL="http://127.0.0.1:${PORT}"

# Colors (safe on macOS Terminal / iTerm2)
GREEN='\033[0;32m'
RED='\033[0;31m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m'

PASS_COUNT=0
step()  { printf "\n${CYAN}==> %s${NC}\n" "$1"; }
pass()  { PASS_COUNT=$((PASS_COUNT + 1)); printf "  ${GREEN}✓${NC} %s\n" "$1"; }
warn()  { printf "  ${YELLOW}!${NC} %s\n" "$1"; }
fail()  { printf "  ${RED}✗ %s${NC}\n" "$1" >&2; exit 1; }

assert_contains() {
    local label="$1"
    local haystack="$2"
    local needle="$3"
    if printf '%s' "$haystack" | grep -Fq "$needle"; then
        pass "$label"
    else
        printf "  ${RED}Expected to find:${NC} %s\n" "$needle" >&2
        printf "  ${RED}Actual output:${NC}    %s\n" "$haystack" >&2
        fail "$label"
    fi
}

cleanup() {
    if [ "$KEEP" -eq 1 ]; then
        printf "\n${YELLOW}Keeping container '%s' running at %s (volume: %s)${NC}\n" "$CONTAINER_NAME" "$BASE_URL" "$VOLUME_NAME"
        printf "Stop later with: docker rm -f %s && docker volume rm %s\n" "$CONTAINER_NAME" "$VOLUME_NAME"
        return
    fi
    docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
    docker volume rm "$VOLUME_NAME" >/dev/null 2>&1 || true
}
trap cleanup EXIT

wait_for_health() {
    local i
    for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30; do
        if curl -fsS "${BASE_URL}/api/health" >/dev/null 2>&1; then
            return 0
        fi
        sleep 1
    done
    echo "--- Container logs ---" >&2
    docker logs "$CONTAINER_NAME" >&2 || true
    fail "Server did not become healthy at ${BASE_URL}/api/health within 30s"
}

step "1. Checking Docker / Colima environment"
if ! command -v docker >/dev/null 2>&1; then
    fail "docker CLI not found in PATH"
fi

if ! docker info >/dev/null 2>&1; then
    if command -v colima >/dev/null 2>&1; then
        warn "Docker daemon not reachable; check Colima with: colima status"
        warn "Start native ARM64 Colima with: colima start --arch aarch64 --cpu 4 --memory 4"
    fi
    fail "Cannot reach Docker daemon"
fi

if command -v colima >/dev/null 2>&1; then
    colima status 2>&1 | sed 's/^/  /' || true
fi

DAEMON_ARCH="$(docker info --format '{{.Architecture}}')"
pass "Docker daemon reachable (arch: ${DAEMON_ARCH})"
if [ "$DAEMON_ARCH" != "$EXPECT_UNAME" ] && [ "$DAEMON_ARCH" != "$EXPECT_ARCH" ]; then
    warn "Docker daemon arch is '${DAEMON_ARCH}', expected '${EXPECT_UNAME}' (running under emulation or different arch)"
fi

step "2. Preparing linux/${EXPECT_ARCH} image (${MODE}: ${IMAGE})"
if [ "$MODE" = "build" ]; then
    SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
    REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
    docker build --platform "linux/${EXPECT_ARCH}" -t "$IMAGE" "$REPO_ROOT"
else
    docker pull --platform "linux/${EXPECT_ARCH}" "$IMAGE"
fi

IMG_ARCH="$(docker image inspect "$IMAGE" --format '{{.Architecture}}')"
assert_contains "Image architecture is ${EXPECT_ARCH}" "$IMG_ARCH" "$EXPECT_ARCH"

step "3. Verifying container CPU & ELF binary architecture"
ARCH_INFO="$(docker run --rm --platform "linux/${EXPECT_ARCH}" --entrypoint bash "$IMAGE" -c 'echo "uname=$(uname -m)"; file /app/ydb-server; file /opt/yottadb/current/libyottadb.so')"
printf '%s\n' "$ARCH_INFO" | sed 's/^/  /'
assert_contains "Container uname -m is ${EXPECT_UNAME}" "$ARCH_INFO" "uname=${EXPECT_UNAME}"
if [ "$EXPECT_ARCH" = "arm64" ]; then
    assert_contains "Go binary is native ARM aarch64 ELF" "$ARCH_INFO" "ARM aarch64"
else
    assert_contains "Go binary is native x86-64 ELF" "$ARCH_INFO" "x86-64"
fi

step "4. Running in-container MUMPS (^ZTEST) + Go test suite"
# Note: using default container FS (native Linux ext4 inside Colima VM, avoiding VirtioFS mmap quirks)
TEST_OUT="$(docker run --rm --platform "linux/${EXPECT_ARCH}" "$IMAGE" test)"
printf '%s\n' "$TEST_OUT" | tail -n 20 | sed 's/^/  /'
assert_contains "Native MUMPS ^ZTEST passed with 0 failures" "$TEST_OUT" "FAIL:0"
assert_contains "Go test suite passed" "$TEST_OUT" "ok  	github.com/jordan-lenchitz/ydb-silly-server"

step "5. Starting live server with persistent Docker volume (${VOLUME_NAME})"
# Named Docker volume stays on Colima's native Linux filesystem for optimal YottaDB shared-memory/mmap safety
docker volume create "$VOLUME_NAME" >/dev/null
docker run -d \
    --name "$CONTAINER_NAME" \
    --hostname ydb-silly-server \
    --platform "linux/${EXPECT_ARCH}" \
    -p "${PORT}:8080" \
    -v "${VOLUME_NAME}:/data" \
    "$IMAGE" >/dev/null

wait_for_health
pass "Server healthy on ${BASE_URL}"

step "6. Exercising HTTP API & YottaDB native MUMPS call-ins"

HEALTH="$(curl -fsS "${BASE_URL}/api/health")"
assert_contains "GET /api/health connected=true" "$HEALTH" '"connected":true'

INTRINSICS="$(curl -fsS "${BASE_URL}/api/sys/intrinsics")"
assert_contains "GET /api/sys/intrinsics reports YottaDB ${EXPECT_UNAME}" "$INTRINSICS" "${EXPECT_UNAME}"
assert_contains "GET /api/sys/intrinsics reports UTF-8" "$INTRINSICS" '"zchset":"UTF-8"'

SET_RES="$(curl -fsS -X POST "${BASE_URL}/api/global/ARMTEST" \
    -H 'Content-Type: application/json' \
    -d '{"subscripts":["mac","colima"],"value":"hello-from-apple-silicon"}')"
assert_contains "POST /api/global/ARMTEST" "$SET_RES" '"status":"OK"'

GET_RES="$(curl -fsS "${BASE_URL}/api/global/ARMTEST?subs=mac,colima")"
assert_contains "GET /api/global/ARMTEST roundtrip" "$GET_RES" '"value":"hello-from-apple-silicon"'

TREE_RES="$(curl -fsS "${BASE_URL}/api/global/ARMTEST/tree?max=10")"
assert_contains "GET /api/global/ARMTEST/tree walks B*-tree" "$TREE_RES" 'hello-from-apple-silicon'

EXEC_RES="$(curl -fsS -X POST "${BASE_URL}/api/execute" \
    -H 'Content-Type: application/json' \
    -d '{"mCode":"W \"MUMPS_ON_\",$ZYRELEASE"}')"
assert_contains "POST /api/execute runs native XECUTE" "$EXEC_RES" "MUMPS_ON_YottaDB"
assert_contains "POST /api/execute reports ${EXPECT_UNAME}" "$EXEC_RES" "${EXPECT_UNAME}"

VM_PROV="$(curl -fsS -X POST "${BASE_URL}/api/vm/provision" \
    -H 'Content-Type: application/json' \
    -d '{"id":"colima-arm-vm","cores":4,"sync":true}')"
assert_contains "POST /api/vm/provision creates VM" "$VM_PROV" '"instanceId":"colima-arm-vm"'

VM_STAT="$(curl -fsS "${BASE_URL}/api/vm/status?id=colima-arm-vm")"
assert_contains "GET /api/vm/status reports provisioned=true" "$VM_STAT" '"provisioned":true'

VM_DEL="$(curl -fsS -X DELETE "${BASE_URL}/api/vm/colima-arm-vm")"
assert_contains "DELETE /api/vm/colima-arm-vm cleans up" "$VM_DEL" '"status":"OK"'

FORGE_RES="$(curl -fsS -X POST "${BASE_URL}/api/forge")"
assert_contains "POST /api/forge dynamically compiles ZEMOJI.m to native ${EXPECT_UNAME}" "$FORGE_RES" '"status":"FORGED"'

DECODE_ASCII="$(curl -fsS "${BASE_URL}/api/decode/A")"
assert_contains "GET /api/decode/A invokes forged ZEMOJI" "$DECODE_ASCII" 'ASCII_HEX_41'

DECODE_EMOJI="$(curl -fsS "${BASE_URL}/api/decode/%F0%9F%98%80")"
assert_contains "GET /api/decode/😀 decodes UTF-8 emoji" "$DECODE_EMOJI" 'EMOJI_HEX_1F600'

CAT_ADD="$(curl -fsS -X POST "${BASE_URL}/api/cats" \
    -H 'Content-Type: application/json' \
    -d '{"name":"APPLE_CAT","mood":"ARM64_PURR"}')"
assert_contains "POST /api/cats adds record" "$CAT_ADD" '"status":"OK"'

CATS_LIST="$(curl -fsS "${BASE_URL}/api/cats")"
assert_contains "GET /api/cats returns persisted record" "$CATS_LIST" 'APPLE_CAT'

MTEST_RES="$(curl -fsS "${BASE_URL}/api/mumps/test")"
assert_contains "GET /api/mumps/test passes inside live server" "$MTEST_RES" '"status":"PASS"'

step "7. Verifying YottaDB global persistence across container restart"
# Graceful stop (SIGTERM) so YottaDB flushes shared-memory buffers cleanly to yottadb.dat
docker stop "$CONTAINER_NAME" >/dev/null
docker rm "$CONTAINER_NAME" >/dev/null
docker run -d \
    --name "$CONTAINER_NAME" \
    --hostname ydb-silly-server \
    --platform "linux/${EXPECT_ARCH}" \
    -p "${PORT}:8080" \
    -v "${VOLUME_NAME}:/data" \
    "$IMAGE" >/dev/null

wait_for_health
PERSIST_RES="$(curl -fsS "${BASE_URL}/api/global/ARMTEST?subs=mac,colima")"
assert_contains "Global ^ARMTEST survived container restart" "$PERSIST_RES" '"value":"hello-from-apple-silicon"'

# Clean up test global before finishing
curl -fsS -X DELETE "${BASE_URL}/api/global/ARMTEST" >/dev/null

printf "\n${GREEN}All %d checks passed on linux/%s (%s)!${NC}\n" "$PASS_COUNT" "$EXPECT_ARCH" "$IMAGE"
