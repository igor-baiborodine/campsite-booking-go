#!/bin/bash
set -e
echo "Creating test campgrounds database..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "test_campgrounds" <<-EOSQL
  CREATE EXTENSION IF NOT EXISTS moddatetime;
EOSQL
