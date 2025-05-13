-- name: CreateSong :one
insert into songs (
    song_id,
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
    delete_date,
    updated_at,
    created_at
)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, now(), now())
returning *;

-- name: GetAllSongs :many
select
    songs.song_id,
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
                        'beatmap_id', beatmaps.beatmap_id,
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
            where beatmaps.song_id = songs.song_id
        ),
        '[]'
    ) as beatmaps
from songs
order by songs.release_date desc, songs.sort desc;

-- name: GetSongByID :one
select
    songs.song_id,
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
                        'beatmapID', beatmaps.beatmap_id,
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
            where beatmaps.song_id = songs.song_id
        ),
        '[]'
    ) as beatmaps
from songs
where songs.song_id = $1;

-- name: GetSongsByTitle :many
select *
from songs
where title = $1;

-- name: GetSongByTitleAndArtist :one
select *
from songs
where title = $1 and artist = $2;

-- name: UpdateSong :one
update songs
set
    title = $1,
    artist = $2,
    genre = $3,
    bpm = $4,
    image_url = $5,
    version = $6,
    is_utage = $7,
    is_available = $8,
    is_new = $9,
    sort = $10,
    release_date = $11,
    delete_date = $12,
    updated_at = now()
where song_id = $13
returning *;
