-- name: CreateBeatmap :one
insert into beatmaps (
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, now(), now())
returning
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at;

-- name: GetAllBeatmaps :many
select
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at
from beatmaps;

-- name: GetBeatmapByBeatmapID :one
select
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at
from beatmaps
where beatmap_id = $1;

-- name: GetBeatmapsBySongID :many
select
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at
from beatmaps
where song_id = $1;

-- name: GetBeatmapBySongIDDifficultyAndType :one
select
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at
from beatmaps
where song_id = $1 and difficulty = $2 and type = $3;

-- name: UpdateBeatmap :one
update beatmaps
set
    song_id = $2,
    difficulty = $3,
    level = $4,
    internal_level = $5,
    type = $6,
    total_notes = $7,
    tap = $8,
    hold = $9,
    slide = $10,
    touch = $11,
    break = $12,
    note_designer = $13,
    max_dx_score = $14,
    updated_at = now()
where beatmap_id = $1
returning
    beatmap_id,
    song_id,
    difficulty,
    level,
    internal_level,
    type,
    total_notes,
    tap,
    hold,
    slide,
    touch,
    break,
    note_designer,
    max_dx_score,
    updated_at,
    created_at;
