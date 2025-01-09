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
    is_valid,
    updated_at,
    created_at
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, now(), now())
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
    is_valid,
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
    is_valid,
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
    is_valid,
    updated_at,
    created_at
from beatmaps
where beatmap_id = $1;

-- name: GetBeatmapBySongID :one
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
    is_valid,
    updated_at,
    created_at
from beatmaps
where song_id = $1;
