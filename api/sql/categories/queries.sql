-- sql/queries.sql

-- ===================== CATEGORIES =====================

-- name: GetCategoryByID :one
SELECT
    c.id, c.name, c.type, c.color, c.user_id,
    c.version, c.created_at, c.updated_at,
    u.id                  AS usuario_id,
    u.nome                AS usuario_nome,
    u.telefone            AS usuario_telefone,
    u.email               AS usuario_email,
    u.is_ativo            AS usuario_is_ativo,
    u.password_hash       AS usuario_password_hash,
    u.version             AS usuario_version,
    u.deleted             AS usuario_deleted,
    u.created_at          AS usuario_created_at,
    u.created_by          AS usuario_created_by,
    u.updated_at          AS usuario_updated_at,
    u.updated_by          AS usuario_updated_by
FROM categories c
LEFT JOIN usuarios u ON c.user_id = u.id AND u.deleted = false
WHERE c.id = $1 AND c.deleted = false;

-- name: ListCategoriesByUser :many
SELECT
    c.id, c.name, c.type, c.color, c.user_id,
    c.version, c.created_at, c.updated_at,
    u.id                  AS usuario_id,
    u.nome                AS usuario_nome,
    u.telefone            AS usuario_telefone,
    u.email               AS usuario_email,
    u.is_ativo            AS usuario_is_ativo,
    u.password_hash       AS usuario_password_hash,
    u.version             AS usuario_version,
    u.deleted             AS usuario_deleted,
    u.created_at          AS usuario_created_at,
    u.created_by          AS usuario_created_by,
    u.updated_at          AS usuario_updated_at,
    u.updated_by          AS usuario_updated_by
FROM categories c
LEFT JOIN usuarios u ON c.user_id = u.id AND u.deleted = false
WHERE c.user_id = $1 AND c.deleted = false
ORDER BY c.name;

-- name: InsertCategory :one
INSERT INTO categories (name, type, color, user_id, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, version, created_at;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, type = $3, color = $4, version = version + 1, updated_at = now(), updated_by = $5
WHERE id = $1 AND user_id = $6 AND deleted = false
RETURNING id, version, updated_at;

-- name: DeleteCategory :exec
UPDATE categories SET deleted = true, version = version + 1, updated_at = now()
WHERE id = $1 AND user_id = $2 AND deleted = false;



