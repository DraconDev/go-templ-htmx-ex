-- name: CreateUser :one
INSERT INTO users (auth_id, email, name, picture, is_admin)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByAuthID :one
SELECT * FROM users WHERE auth_id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: GetAdminUsers :many
SELECT * FROM users WHERE is_admin = true ORDER BY created_at DESC;

-- name: UpdateUserAdminStatus :exec
UPDATE users SET is_admin = $2 WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users 
SET name = COALESCE($2, name),
    picture = COALESCE($3, picture),
    updated_at = NOW()
WHERE id = $1;