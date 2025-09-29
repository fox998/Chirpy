-- name: PostChirp :one
INSERT INTO chirps (body, user_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: ListChirpsAll :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: ListChirpsByUser :many
SELECT * FROM chirps
WHERE user_id = $1::uuid
ORDER BY created_at;



-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps
WHERE id = $1 AND user_id = $2;
