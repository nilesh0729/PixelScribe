-- name: CreateSetting :one
INSERT INTO settings (
    user_id, default_voice, default_speed, 
    highlight_color_grammar, highlight_color_spelling, highlight_color_case, 
    created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
)
RETURNING id, user_id, default_voice, default_speed, highlight_color_grammar, highlight_color_spelling, highlight_color_case, created_at, updated_at;

-- name: GetSettingByID :one
SELECT id, user_id, default_voice, default_speed, highlight_color_grammar, highlight_color_spelling, highlight_color_case, created_at, updated_at
FROM settings
WHERE id = $1;

-- name: GetSettingByUserID :one
SELECT id, user_id, default_voice, default_speed, highlight_color_grammar, highlight_color_spelling, highlight_color_case, created_at, updated_at
FROM settings
WHERE user_id = $1;

-- name: UpdateSetting :one
UPDATE settings
SET 
    default_voice = $2,
    default_speed = $3,
    highlight_color_grammar = $4,
    highlight_color_spelling = $5,
    highlight_color_case = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, default_voice, default_speed, highlight_color_grammar, highlight_color_spelling, highlight_color_case, created_at, updated_at;

-- name: DeleteSetting :exec
DELETE FROM settings
WHERE id = $1;
