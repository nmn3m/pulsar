#!/bin/bash

# Pulsar Demo Data Seeder
# This script seeds the database with demo data for presentations

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

echo "=================================="
echo "Pulsar Demo Data Seeder"
echo "=================================="
echo ""

# Check if required environment variables are set
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is required"
    echo ""
    echo "Example:"
    echo "  export DATABASE_URL='postgres://user:pass@localhost:5432/pulsar?sslmode=disable'"
    exit 1
fi

if [ -z "$JWT_SECRET" ]; then
    echo "Warning: JWT_SECRET not set, using default for demo"
    export JWT_SECRET="demo-secret-key-at-least-32-chars"
fi

if [ -z "$JWT_REFRESH_SECRET" ]; then
    echo "Warning: JWT_REFRESH_SECRET not set, using default for demo"
    export JWT_REFRESH_SECRET="demo-refresh-secret-at-least-32c"
fi

# Build and run the seed command
echo "Building seed command..."
cd "$BACKEND_DIR"
go build -o bin/seed ./cmd/seed

echo "Running seed..."
./bin/seed

echo ""
echo "=================================="
echo "Demo data seeded successfully!"
echo "=================================="
