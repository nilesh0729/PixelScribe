-- name: CreateTextDictations :one
INSERT INTO dictations (
  user_id,
  title,
  type,
  content, 
  language,
  created_at, 
  updated_at

) VALUES (
  $1, $2, 'text', $3, $4, $5, $6
)
RETURNING *;

-- name: CreateAudioDictations :one
INSERT INTO dictations (
  user_id,
  title,
  type,
  audio_url,
  language,
  created_at,
  updated_at
) VALUES (
  $1, $2, 'audio', $3, $4, $5, $6
)
RETURNING *;

-- name: GetDictationsByTitle :one
SELECT * FROM dictations
WHERE title = $1 LIMIT 1;

-- name: GetDictation :one
SELECT * FROM dictations
WHERE id = $1 LIMIT 1;

-- name: ListDictationsByUser :many
SELECT * FROM dictations
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListTextDictations :many
SELECT * FROM dictations
WHERE user_id = $1
    AND type = 'text'
ORDER BY created_at DESC;

-- name: ListAudioDictations :many
SELECT * FROM dictations
WHERE user_id = $1
    AND type = 'audio'
ORDER BY created_at DESC;

-- name: UpdateDictation :one
UPDATE dictations
SET
    title = COALESCE(sqlc.narg('title'), title),
    content = COALESCE(sqlc.narg('content'), content),
    audio_url = COALESCE(sqlc.narg('audio_url'), audio_url),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
  AND user_id = sqlc.arg('user_id')
RETURNING *;


-- name: DeleteDictations :exec
DELETE FROM dictations
WHERE title = $1;