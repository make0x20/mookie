-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES (?, ?, ?)
RETURNING id, username, email, password, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ? LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
