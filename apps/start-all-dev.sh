#!/bin/bash

set -e

AUTH_DIR="./auth"
CLIENTS_DIR="./clients"

# Function to clean up on Ctrl+C
cleanup() {
  echo -e "\nStopping all microservices..."
  kill $AUTH_PID $CLIENTS_PID 2>/dev/null || true
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

echo "Auth PID: $AUTH_PID"
echo "Clients PID: $CLIENTS_PID"
echo "Press Ctrl+C to stop both microservices."

# Wait for both
wait $AUTH_PID $CLIENTS_PID
