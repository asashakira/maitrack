-- name: CreateScore :one
insert into scores (
    id,
    beatmap_id,
    song_id,
    user_uuid,
    accuracy,
    max_combo,
    dx_score,
    tap_critical,
    tap_perfect,
    tap_great,
    tap_good,
    tap_miss,
    hold_critical,
    hold_perfect,
    hold_great,
    hold_good,
    hold_miss,
    slide_critical,
    slide_perfect,
    slide_great,
    slide_good,
    slide_miss,
    touch_critical,
    touch_perfect,
    touch_great,
    touch_good,
    touch_miss,
    break_critical,
    break_perfect,
    break_great,
    break_good,
    break_miss,
    fast,
    late,
    played_at
)
values (
    $1, $2, $3, $4, $5, $6, $7,
    $8, $9, $10, $11, $12,
    $13, $14, $15, $16, $17,
    $18, $19, $20, $21, $22,
    $23, $24, $25, $26, $27,
    $28, $29, $30, $31, $32,
    $33, $34, $35
)
returning *;


-- name: GetScoresByUserID :many
select
    scores.id,
    scores.beatmap_id,
    scores.song_id,
    scores.user_uuid,
    scores.accuracy,
    scores.max_combo,
    scores.dx_score,
    scores.tap_critical,
    scores.tap_perfect,
    scores.tap_great,
    scores.tap_good,
    scores.tap_miss,
    scores.hold_critical,
    scores.hold_perfect,
    scores.hold_great,
    scores.hold_good,
    scores.hold_miss,
    scores.slide_critical,
    scores.slide_perfect,
    scores.slide_great,
    scores.slide_good,
    scores.slide_miss,
    scores.touch_critical,
    scores.touch_perfect,
    scores.touch_great,
    scores.touch_good,
    scores.touch_miss,
    scores.break_critical,
    scores.break_perfect,
    scores.break_great,
    scores.break_good,
    scores.break_miss,
    scores.fast,
    scores.late,
    scores.played_at,
    songs.title,
    songs.artist,
    songs.genre,
    songs.image_url,
    songs.version,
    beatmaps.difficulty,
    beatmaps.level,
    beatmaps.internal_level,
    beatmaps.type
from scores
inner join songs on scores.song_id = songs.id
inner join beatmaps on scores.beatmap_id = beatmaps.id
inner join users on scores.user_uuid = users.id
where users.user_id = $1
order by scores.played_at desc
limit $2 offset $3;
