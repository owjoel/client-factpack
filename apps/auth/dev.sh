#!/bin/bash
clear

# Check if the .env file exists
if [ ! -f .env ]; then
  echo ".env file not found."
  exit 1
fi

set -a
source .env
set +a

go run ./cmd