-- name: CreateSong :one
insert into songs (
    id,
    alt_key,
    title,
    artist,
    genre,
    bpm,
    image_url,
    version,
    is_utage,
    is_available,
    is_new,
    sort,
    release_date,
    delete_date
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
returning *;

-- name: GetAllSongs :many
select
    songs.id,
    songs.title,
    songs.artist,
    songs.genre,
    songs.bpm,
    songs.image_url,
    songs.version,
    songs.is_utage,
    songs.is_available,
    songs.is_new,
    songs.release_date,
    songs.delete_date,
    coalesce(
        (
            select
                json_agg(
                    jsonb_build_object(
                        'beatmap_id', beatmaps.id,
                        'difficulty', beatmaps.difficulty,
                        'level', beatmaps.level,
                        'internal_level', beatmaps.internal_level,
                        'type', beatmaps.type,
                        'total_notes', beatmaps.total_notes,
                        'tap', beatmaps.tap,
                        'hold', beatmaps.hold,
                        'slide', beatmaps.slide,
                        'touch', beatmaps.touch,
                        'break', beatmaps.break,
                        'note_designer', beatmaps.note_designer,
                        'max_dx_score', beatmaps.max_dx_score
                    )
                )
            from beatmaps
            where beatmaps.song_id = songs.id
        ),
        '[]'
    ) as beatmaps
from songs
order by songs.release_date desc, songs.sort desc;

-- name: GetSongByID :one
select
    songs.id,
    songs.title,
    songs.artist,
    songs.genre,
    songs.bpm,
    songs.image_url,
    songs.version,
    songs.is_utage,
    songs.is_available,
    songs.is_new,
    songs.release_date,
    songs.delete_date,
    coalesce(
        (
            select
                json_agg(
                    jsonb_build_object(
                        'beatmapID', beatmaps.id,
                        'difficulty', beatmaps.difficulty,
                        'level', beatmaps.level,
                        'internalLevel', beatmaps.internal_level,
                        'type', beatmaps.type,
                        'totalNotes', beatmaps.total_notes,
                        'tap', beatmaps.tap,
                        'hold', beatmaps.hold,
                        'slide', beatmaps.slide,
                        'touch', beatmaps.touch,
                        'break', beatmaps.break,
                        'noteDesigner', beatmaps.note_designer,
                        'maxDxScore', beatmaps.max_dx_score
                    )
                )
            from beatmaps
            where beatmaps.song_id = songs.id
        ),
        '[]'
    ) as beatmaps
from songs
where songs.id = $1;

-- name: GetSongsByTitle :many
select *
from songs
where title = $1;

-- name: GetSongByTitleAndArtist :one
select *
from songs
where title = $1 and artist = $2;

-- name: GetSongByAltKey :one
select *
from songs
where alt_key = $1;

-- name: UpdateSong :one
update songs
set
    alt_key = $1,
    title = $2,
    artist = $3,
    genre = $4,
    bpm = $5,
    image_url = $6,
    version = $7,
    is_utage = $8,
    is_available = $9,
    is_new = $10,
    sort = $11,
    release_date = $12,
    delete_date = $13,
    updated_at = now()
where id = $14
returning *;
