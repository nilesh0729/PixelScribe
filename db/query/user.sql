-- name: CreateUsers :one
INSERT INTO users (
  name,
  username, 
  email, 
  password_hash
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: GetUsers :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name
LIMIT $1
OFFSET $2;

-- name: UpdateUsers :one
UPDATE users
  set name = $2,
  email = $3,
  password_hash = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users
WHERE username = $1;
