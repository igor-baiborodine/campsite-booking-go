#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE campgrounds;
  CREATE USER campgrounds_user WITH ENCRYPTED PASSWORD 'campgrounds_pass';
  GRANT USAGE ON SCHEMA public TO campgrounds_user;
  GRANT CREATE, CONNECT ON DATABASE campgrounds TO campgrounds_user;
EOSQL

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "campgrounds" <<-EOSQL
  CREATE SCHEMA campgrounds;
  GRANT CREATE, USAGE ON SCHEMA campgrounds TO campgrounds_user;
EOSQL
