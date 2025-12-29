-- name: CreateAttempts :one
INSERT INTO attempts (
  user_id, 
  dictation_id, 
  typed_text, 
  attempt_no, 
  total_words, 
  correct_words, 
  grammatical_errors, 
  spelling_errors, 
  case_errors, 
  accuracy, 
  comparison_data, 
  time_spent,
  created_at
) VALUES (
  $1, $2, $3, 
  COALESCE((
    SELECT MAX(attempt_no) + 1
    FROM attempts
    WHERE user_id = $1 AND dictation_id = $2
  ), 1), 
  $4, $5, $6, $7, $8, $9, $10, $11, NOW()
)
RETURNING *;

-- name: GetAttemptById :one
SELECT * FROM attempts
WHERE id = $1 LIMIT 1;

-- name: ListAttemptsByUser :many
SELECT * FROM attempts
WHERE user_id = $1
ORDER by created_at DESC;

-- name: ListAttemptsByDictation :many
SELECT * FROM attempts
WHERE dictation_id = $1
ORDER BY created_at DESC;

-- name: GetLatestAttempt :one
SELECT * FROM attempts
WHERE user_id = $1 AND dictation_id = $2
ORDER BY created_at DESC;

-- name: CountAttemptsByDictation :one
SELECT COUNT(*) FROM attempts
WHERE user_id = $1 AND dictation_id = $2;

-- name: DeleteAttempt :exec
DELETE FROM attempts
WHERE id = $1;

-- name: DeleteAttemptsByDictation :exec
DELETE FROM attempts
WHERE dictation_id = $1 AND user_id = $2;

-- name: UpdateAttemptAccuracy :one
UPDATE attempts
SET
  accuracy = $2,
  correct_words = $3,
  grammatical_errors = $4,
  spelling_errors = $5,
  case_errors = $6,
  comparison_data = $7,
  time_spent = $8
WHERE id = $1
RETURNING *;