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

-- name: UpdateUser :one
UPDATE users
SET
    email = $2,
    password_hash = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
