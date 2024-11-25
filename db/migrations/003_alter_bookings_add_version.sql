-- +goose Up
ALTER TABLE bookings ADD COLUMN version INT DEFAULT 1;

-- +goose Down
ALTER TABLE bookings DROP COLUMN IF EXISTS version;
