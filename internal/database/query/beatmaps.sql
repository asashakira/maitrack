-- name: CreateBeatmap :one
insert into beatmaps (
    id,
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
    max_dx_score
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
returning *;


-- name: GetAllBeatmaps :many
select
    beatmaps.id,
    beatmaps.song_id,
    beatmaps.difficulty,
    beatmaps.level,
    beatmaps.internal_level,
    beatmaps.type,
    beatmaps.total_notes,
    beatmaps.tap,
    beatmaps.hold,
    beatmaps.slide,
    beatmaps.touch,
    beatmaps.break,
    beatmaps.note_designer,
    beatmaps.max_dx_score,
    songs.title,
    songs.artist,
    songs.genre,
    songs.bpm,
    songs.image_url,
    songs.version
from beatmaps
inner join songs on beatmaps.song_id = songs.id;


-- name: GetBeatmapByBeatmapID :one
select
    beatmaps.id,
    beatmaps.song_id,
    beatmaps.difficulty,
    beatmaps.level,
    beatmaps.internal_level,
    beatmaps.type,
    beatmaps.total_notes,
    beatmaps.tap,
    beatmaps.hold,
    beatmaps.slide,
    beatmaps.touch,
    beatmaps.break,
    beatmaps.note_designer,
    beatmaps.max_dx_score,
    songs.title,
    songs.artist,
    songs.genre,
    songs.bpm,
    songs.image_url,
    songs.version
from beatmaps
inner join songs on beatmaps.song_id = songs.id
where beatmaps.id = $1;


-- name: GetBeatmapsBySongID :many
select
    beatmaps.id,
    beatmaps.song_id,
    beatmaps.difficulty,
    beatmaps.level,
    beatmaps.internal_level,
    beatmaps.type,
    beatmaps.total_notes,
    beatmaps.tap,
    beatmaps.hold,
    beatmaps.slide,
    beatmaps.touch,
    beatmaps.break,
    beatmaps.note_designer,
    beatmaps.max_dx_score,
    songs.title,
    songs.artist,
    songs.genre,
    songs.bpm,
    songs.image_url,
    songs.version
from beatmaps
inner join songs on beatmaps.song_id = songs.song_id
where beatmaps.song_id = $1;

-- name: GetBeatmapBySongIDDifficultyAndType :one
select *
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
where id = $1
returning *;
