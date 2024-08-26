#!/bin/bash
set -e
echo "Creating campgrounds database..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE campgrounds;
  CREATE USER campgrounds_user WITH ENCRYPTED PASSWORD 'campgrounds_pass';
  ALTER DATABASE campgrounds OWNER TO campgrounds_user;
  GRANT USAGE ON SCHEMA public TO campgrounds_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "campgrounds" <<-EOSQL
  CREATE EXTENSION IF NOT EXISTS moddatetime;
EOSQL
