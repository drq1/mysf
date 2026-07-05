#!/bin/sh
set -e

REDIS_BIND="${REDIS_BIND:-127.0.0.1}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_DATA_DIR="${REDIS_DATA_DIR:-/app/data/redis}"
REDIS_APPENDONLY="${REDIS_APPENDONLY:-yes}"

cleanup() {
    if [ -n "${app_pid:-}" ] && kill -0 "$app_pid" 2>/dev/null; then
        kill "$app_pid" 2>/dev/null || true
    fi
    if [ -n "${redis_pid:-}" ] && kill -0 "$redis_pid" 2>/dev/null; then
        kill "$redis_pid" 2>/dev/null || true
    fi
}

# Fix data directory permissions when running as root.
if [ "$(id -u)" = "0" ]; then
    mkdir -p /app/data /app/data/redis
    chown -R sub2api:sub2api /app/data 2>/dev/null || true
    exec su-exec sub2api "$0" "$@"
fi

# Compatibility: if the first arg looks like a flag, prepend the default binary.
if [ "${1#-}" != "$1" ]; then
    set -- /app/sub2api "$@"
fi

# Preserve the normal "docker run image sh" behavior.
if [ "$1" != "/app/sub2api" ]; then
    exec "$@"
fi

mkdir -p "$REDIS_DATA_DIR"

redis-server \
    --bind "$REDIS_BIND" \
    --port "$REDIS_PORT" \
    --dir "$REDIS_DATA_DIR" \
    --loglevel "${REDIS_LOGLEVEL:-warning}" \
    --appendonly "$REDIS_APPENDONLY" \
    --save 60 1000 \
    --protected-mode yes \
    >/dev/null 2>&1 &
redis_pid=$!

trap cleanup INT TERM EXIT

redis_ready=0
i=0
while [ $i -lt 60 ]; do
    if redis-cli -h "$REDIS_BIND" -p "$REDIS_PORT" ping >/dev/null 2>&1; then
        redis_ready=1
        break
    fi
    i=$((i + 1))
    sleep 1
done

if [ "$redis_ready" -ne 1 ]; then
    echo "Redis did not become ready in time" >&2
    exit 1
fi

"$@" >/dev/null 2>&1 &
app_pid=$!

wait "$app_pid"
app_status=$?

cleanup
trap - INT TERM EXIT
exit "$app_status"
