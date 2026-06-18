-- sql/queries.sql

-- ===================== USUÁRIOS =====================

-- name: GetUserByID :one
SELECT id, nome, telefone, email, is_ativo, password_hash, version, deleted, created_at, created_by, updated_at, updated_by
FROM usuarios
WHERE id = $1 AND deleted = false;

-- name: GetUserByEmail :one
SELECT id, nome, telefone, email, is_ativo, password_hash, version, deleted, created_at, created_by, updated_at, updated_by
FROM usuarios
WHERE email = $1 AND deleted = false;

-- name: InsertUser :one
INSERT INTO usuarios (nome, telefone, email, is_ativo, password_hash, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, version, created_at;

-- name: UpdateUser :one
UPDATE usuarios
SET nome = $2, telefone = $3, email = $4, is_ativo = $5, password_hash = $6,
    version = version + 1, updated_at = now(), updated_by = $7
WHERE id = $1 AND deleted = false
RETURNING id, version, updated_at;

-- name: DeleteUser :exec
UPDATE usuarios SET deleted = true, version = version + 1, updated_at = now()
WHERE id = $1 AND deleted = false;



