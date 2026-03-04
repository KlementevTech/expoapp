-- name: CreatePart :exec
INSERT INTO expo_parts (id, name, version, created_at)
VALUES (:id, :name, :version, :created_at);

-- name: GetPart :one
SELECT id, name, version, created_at
FROM expo_parts
WHERE id = :id AND deleted_at IS NULL
LIMIT 1;

-- name: ListParts :many
SELECT id, name, version, created_at, updated_at
FROM expo_parts
WHERE deleted_at IS NULL
ORDER BY created_at DESC, id DESC;

-- name: UpdatePart :execresult
UPDATE expo_parts
SET name = :name,
    version = :version,
    updated_at = :updated_at
WHERE id = :id AND version = :old_version AND deleted_at IS NULL;

-- name: DeletePart :execresult
UPDATE expo_parts
SET version = version + 1,
    updated_at = :updated_at,
    deleted_at = :deleted_at
WHERE id = :id AND version = :old_version AND deleted_at IS NULL;

-- name: ExistsPart :one
SELECT CAST(EXISTS(
    SELECT 1 FROM expo_parts
    WHERE id = :id AND deleted_at IS NULL
) AS BOOLEAN);

-- name: CleanParts :exec
DELETE FROM expo_parts;