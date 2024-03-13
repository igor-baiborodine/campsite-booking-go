-- +goose Up
CREATE SCHEMA IF NOT EXISTS campgrounds;

SET SEARCH_PATH TO campgrounds, public;

-- +goose Down
DROP SCHEMA IF EXISTS campgrounds CASCADE;
