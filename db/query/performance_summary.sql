-- name: CreatePerformanceSummary :one
INSERT INTO performance_summary (
    user_id, dictation_id, total_attempts, best_accuracy, average_accuracy, average_time, last_attempt_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetPerformanceSummaryByID :one
SELECT *
FROM performance_summary
WHERE id = $1;

-- name: ListPerformanceSummaryByUser :many
SELECT *
FROM performance_summary
WHERE user_id = $1
ORDER BY last_attempt_at DESC;

-- name: GetPerformanceSummaryByUserAndDictation :one
SELECT *
FROM performance_summary
WHERE user_id = $1 AND dictation_id = $2;

-- name: UpdatePerformanceSummary :one
UPDATE performance_summary
SET
    total_attempts = $1,
    best_accuracy = $2,
    average_accuracy = $3,
    average_time = $4,
    last_attempt_at = $5
WHERE id = $6
RETURNING *;

-- name: DeletePerformanceSummary :exec
DELETE FROM performance_summary
WHERE id = $1;

-- name: RecentAttemptsByUser :many
SELECT *
FROM performance_summary
WHERE user_id = $1
ORDER BY last_attempt_at DESC
LIMIT $2;

-- name: UserAggregatePerformance :many
SELECT 
    user_id,
    AVG(average_accuracy)::float AS overall_avg_accuracy,
    AVG(average_time)::float AS overall_avg_time
FROM performance_summary
GROUP BY user_id;


