#!/bin/bash

set -e

AUTH_DIR="./apps/auth"
CLIENTS_DIR="./apps/clients"
PREFECT_DIR="./prefect"

# Function to clean up on Ctrl+C
cleanup() {
  echo -e "\nStopping all microservices..."
  kill $AUTH_PID $CLIENTS_PID $PREFECT_PID 2>/dev/null || true
  exit 0
}

trap cleanup SIGINT SIGTERM

# Start auth service
cd "$AUTH_DIR"
chmod +x ./dev.sh
./dev.sh 2>&1 | sed 's/^/[AUTH] /' &
AUTH_PID=$!
cd - > /dev/null

# Start clients service
cd "$CLIENTS_DIR"
chmod +x ./dev.sh
./dev.sh 2>&1 | sed 's/^/[CLIENTS] /' &
CLIENTS_PID=$!
cd - > /dev/null

# Start Prefect worker
cd "$PREFECT_DIR"

# Load environment variables from .env file
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
else
  echo "[PREFECT] .env file not found in $PREFECT_DIR. Exiting."
  exit 1
fi

# Install Python dependencies if not already installed
if ! pip list | grep -q "prefect"; then
  echo "[PREFECT] Installing dependencies..."
  pip install -r requirements.txt
fi

# Start the Prefect worker
echo "[PREFECT] Starting Prefect worker..."
prefect worker start --pool justin-local 2>&1 | sed 's/^/[PREFECT] /' &
PREFECT_PID=$!
cd - > /dev/null

echo "Auth PID: $AUTH_PID"
echo "Clients PID: $CLIENTS_PID"
echo "Prefect PID: $PREFECT_PID"
echo "Press Ctrl+C to stop all microservices."

# Wait for all
wait $AUTH_PID $CLIENTS_PID $PREFECT_PID
