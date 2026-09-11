-- ============================================================================
-- GESTIÓN DE USUARIOS
-- ============================================================================

-- name: CreateUser :one
INSERT INTO users (
    name
) VALUES (
    $1
)
RETURNING id, name, created_at;

-- name: GetUser :one
SELECT id, name, created_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, created_at
FROM users
ORDER BY id ASC;

-- name: GetRandomUser :one
SELECT id, name, created_at
FROM users
ORDER BY RANDOM()
LIMIT 1;

-- name: UpdateUser :exec
UPDATE users
SET name = $2
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
