-- name: CreateUserData :one
insert into user_data (
    id,
    user_uuid,
    rating,
    season_play_count,
    total_play_count
)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetUserDataByUserUUID :one
select
    id,
    user_uuid,
    rating,
    season_play_count,
    total_play_count,
    created_at
from user_data
where user_uuid = $1
order by created_at desc
limit 1;
