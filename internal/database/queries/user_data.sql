-- name: CreateUserData :one
insert into user_data (
    id,
    user_id,
    game_name,
    tag_line,
    rating,
    season_play_count,
    total_play_count,
    created_at
)
values ($1, $2, $3, $4, $5, $6, $7, now())
returning
    id,
    user_id,
    game_name,
    tag_line,
    rating,
    season_play_count,
    total_play_count,
    created_at;

-- name: GetUserDataByUserID :one
select
    id,
    user_id,
    game_name,
    tag_line,
    rating,
    season_play_count,
    total_play_count,
    created_at
from user_data
where user_id = $1
order by created_at desc
limit 1;

-- name: GetUserDataByMaiID :one
select
    id,
    user_id,
    game_name,
    tag_line,
    rating,
    season_play_count,
    total_play_count,
    created_at
from user_data
where game_name = $1 and tag_line = $2
order by created_at desc
limit 1;
