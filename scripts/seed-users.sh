#!/bin/bash

# Seed users script for Algorithmia
# This script creates test users for all roles and generates a credentials file

set -e

echo "🌱 Seeding test users..."
echo ""

# Load environment variables from .env if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run the seed script
cd cmd/tools/seed-users
go run main.go

# Move the credentials file to the project root
if [ -f test_users_credentials.txt ]; then
    mv test_users_credentials.txt ../../../
    echo ""
    echo "✅ Done! Check test_users_credentials.txt for login credentials"
else
    echo ""
    echo "❌ Credentials file not found"
fi
