-- name: InsertUserData :exec
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
values ($1, $2, $3, $4, $5, $6, $7, now());

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
where user_id = $1;

-- name: DeleteUserDataByUserID :exec
delete from user_data
where user_id = $1;
