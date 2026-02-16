-- +goose Up
CREATE TABLE IF NOT EXISTS device (
    id BIGSERIAL PRIMARY KEY,
    unit_guid UUID NOT NULL,
    n INTEGER,
    mqtt TEXT,
    invid TEXT,
    msg_id TEXT,
    text TEXT,
    context TEXT,
    class INTEGER,
    level TEXT,
    area TEXT,
    addr TEXT,
    block TEXT,
    type TEXT,
    bit TEXT,
    invert_bit TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_device_unit_guid ON device(unit_guid);

CREATE TABLE IF NOT EXISTS processed_file (
    id BIGSERIAL PRIMARY KEY,
    file_name TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS parse_error (
    id BIGSERIAL PRIMARY KEY,
    file_name TEXT NOT NULL,
    line_number INTEGER,
    error_message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down
--DROP TABLE IF EXISTS parse_error;
--DROP TABLE IF EXISTS processed_file;
--DROP TABLE IF EXISTS device;
