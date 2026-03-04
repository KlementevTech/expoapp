-- +goose Up
CREATE TABLE IF NOT EXISTS expo_parts (
    id         BLOB PRIMARY KEY, -- UUIDv7
    name       TEXT NOT NULL,
    version    INTEGER NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    deleted_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_parts_active_id
    ON expo_parts (id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_parts_list
    ON expo_parts (created_at DESC, id DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_parts_list;
DROP INDEX IF EXISTS idx_parts_active_id;
DROP TABLE IF EXISTS expo_parts;
