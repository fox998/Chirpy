-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: ResetUsers :exec
TRUNCATE TABLE users CASCADE;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;
