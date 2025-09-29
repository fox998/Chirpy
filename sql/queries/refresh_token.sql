
-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, expaires_at)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetRefreshTokenById :one
SELECT *
FROM refresh_tokens
WHERE id = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE id = $1;