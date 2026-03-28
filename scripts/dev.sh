#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

BACKEND_PORT=8080
FRONTEND_PORT=3000
DB_PORT=5432

cleanup() {
    echo "Cleaning up..."

    if [ -n "$BACKEND_PID" ]; then
        echo "Stopping backend (PID: $BACKEND_PID)..."
        kill $BACKEND_PID 2>/dev/null || true
    fi

    if [ -n "$FRONTEND_PID" ]; then
        echo "Stopping frontend (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID 2>/dev/null || true
    fi

    if [ "$STOP_DB" = "true" ]; then
        echo "Stopping PostgreSQL..."
        docker stop army-game-db 2>/dev/null || true
    fi

    echo "Cleanup complete"
}

check_port() {
    local port=$1
    if lsof -i:$port >/dev/null 2>&1; then
        return 0
    fi
    return 1
}

kill_port() {
    local port=$1
    echo "Port $port is in use, killing..."

    local pids=$(lsof -t -i:$port 2>/dev/null || true)
    if [ -n "$pids" ]; then
        echo "Killing processes on port $port: $pids"
        echo "$pids" | xargs kill -9 2>/dev/null || true
        sleep 1
    fi

    if check_port $port; then
        echo "Failed to free port $port"
        return 1
    fi

    echo "Port $port is now free"
    return 0
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "Docker not found"
        return 1
    fi

    if ! docker ps &> /dev/null; then
        echo "Docker daemon not running"
        return 1
    fi

    return 0
}

wait_for_backend() {
    echo "Waiting for backend to be ready..."
    local max_attempts=30
    local attempt=0

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:$BACKEND_PORT/health >/dev/null 2>&1; then
            echo "Backend is ready!"
            return 0
        fi
        sleep 1
        attempt=$((attempt + 1))
    done

    echo "Backend failed to start"
    return 1
}

STOP_DB=false

start_postgres() {
    if check_docker; then
        if docker ps --format '{{.Names}}' | grep -q "^army-game-db$"; then
            echo "PostgreSQL container already running"
        elif docker ps -a --format '{{.Names}}' | grep -q "^army-game-db$"; then
            echo "Starting existing PostgreSQL container..."
            docker start army-game-db
        else
            echo "Starting PostgreSQL..."
            docker-compose up -d postgres
        fi

        echo "Waiting for PostgreSQL to be ready..."
        sleep 3
    else
        echo "Docker not available, skipping PostgreSQL"
    fi
}

start_backend() {
    echo "Starting backend on port $BACKEND_PORT..."

    if check_port $BACKEND_PORT; then
        kill_port $BACKEND_PORT
    fi

    cd "$PROJECT_ROOT/backend"
    go run cmd/server/main.go &
    BACKEND_PID=$!

    wait_for_backend
}

start_frontend() {
    echo "Starting frontend on port $FRONTEND_PORT..."

    if check_port $FRONTEND_PORT; then
        kill_port $FRONTEND_PORT
    fi

    cd "$PROJECT_ROOT/frontend"
    npm run dev &
    FRONTEND_PID=$!

    echo "Frontend started (PID: $FRONTEND_PID)"
}

usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start       Start all services (default)"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  status      Check service status"
    echo "  test        Run all tests"
    echo "  clean       Clean up ports and temp files"
}

cmd_start() {
    trap cleanup EXIT

    start_postgres

    start_backend

    start_frontend

    echo ""
    echo "========================================="
    echo "Services started:"
    echo "  Backend:   http://localhost:$BACKEND_PORT"
    echo "  Frontend:  http://localhost:$FRONTEND_PORT"
    echo "  GraphQL:   http://localhost:$BACKEND_PORT/graphql"
    echo "========================================="
    echo ""
    echo "Press Ctrl+C to stop all services"

    wait
}

cmd_stop() {
    echo "Stopping all services..."

    local backend_pid=$(lsof -t -i:$BACKEND_PORT 2>/dev/null || true)
    if [ -n "$backend_pid" ]; then
        echo "Stopping backend (PID: $backend_pid)..."
        kill $backend_pid 2>/dev/null || true
    fi

    local frontend_pid=$(lsof -t -i:$FRONTEND_PORT 2>/dev/null || true)
    if [ -n "$frontend_pid" ]; then
        echo "Stopping frontend (PID: $frontend_pid)..."
        kill $frontend_pid 2>/dev/null || true
    fi

    echo "All services stopped"
}

cmd_status() {
    echo "Service Status:"
    echo ""

    if check_port $BACKEND_PORT; then
        echo "  Backend:   RUNNING (port $BACKEND_PORT)"
    else
        echo "  Backend:   STOPPED (port $BACKEND_PORT)"
    fi

    if check_port $FRONTEND_PORT; then
        echo "  Frontend: RUNNING (port $FRONTEND_PORT)"
    else
        echo "  Frontend: STOPPED (port $FRONTEND_PORT)"
    fi

    if check_docker && docker ps --format '{{.Names}}' | grep -q "^army-game-db$"; then
        echo "  PostgreSQL: RUNNING (port $DB_PORT)"
    else
        echo "  PostgreSQL: STOPPED (port $DB_PORT)"
    fi
}

cmd_clean() {
    echo "Cleaning up..."

    if check_port $BACKEND_PORT; then
        kill_port $BACKEND_PORT
    fi

    if check_port $FRONTEND_PORT; then
        kill_port $FRONTEND_PORT
    fi

    echo "Cleanup complete"
}

cmd_test() {
    echo "Running tests..."

    echo ""
    echo "=== Backend Tests ==="
    cd "$PROJECT_ROOT/backend"
    go test -v ./...

    echo ""
    echo "=== Frontend Tests ==="
    cd "$PROJECT_ROOT/frontend"
    npm test

    echo ""
    echo "All tests complete"
}

COMMAND=${1:-start}

case $COMMAND in
    start)
        STOP_DB=true
        cmd_start
        ;;
    stop)
        cmd_stop
        ;;
    restart)
        cmd_stop
        sleep 2
        STOP_DB=true cmd_start
        ;;
    status)
        cmd_status
        ;;
    clean)
        cmd_clean
        ;;
    test)
        cmd_test
        ;;
    *)
        usage
        exit 1
        ;;
esac
