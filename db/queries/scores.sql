-- name: CreateScore :one
insert into scores (
    score_id,
    beatmap_id,
    song_id,
    user_id,
    accuracy,
    max_combo,
    dx_score,
    played_at,
    created_at
)
values ($1, $2, $3, $4, $5, $6, $7, $8, now())
returning *;

-- name: GetAllScores :many
select
    score_id,
    beatmap_id,
    song_id,
    user_id,
    accuracy,
    max_combo,
    dx_score,
    played_at,
    created_at
from scores;

-- name: GetScoreByUserID :one
select
    score_id,
    beatmap_id,
    song_id,
    user_id,
    accuracy,
    max_combo,
    dx_score,
    played_at,
    created_at
from scores
where user_id = $1;
