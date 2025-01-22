-- name: CreateUser :one
insert into users (
    user_id,
    username,
    password,
    sega_id,
    sega_password,
    game_name,
    tag_line,
    updated_at,
    created_at
)
values ($1, $2, $3, $4, $5, $6, $7, now(), now())
returning user_id, username, password, sega_id, sega_password, game_name, tag_line, updated_at, created_at;

-- name: GetAllUsers :many
select
    users.user_id,
    users.username,
    users.password,
    users.sega_id,
    users.sega_password,
    users.game_name,
    users.tag_line,
    user_data.rating,
    user_data.season_play_count,
    user_data.total_play_count
from users
left join (
    select distinct on (user_data.user_id)
        user_data.user_id,
        user_data.rating,
        user_data.season_play_count,
        user_data.total_play_count,
        user_data.created_at
    from user_data
    order by user_data.user_id asc, user_data.created_at desc
) as user_data
    on users.user_id = user_data.user_id;

-- name: GetUserByID :one
select
    user_id,
    username,
    password,
    sega_id,
    sega_password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
where user_id = $1;

-- name: GetUserByUsername :one
select
    user_id,
    username,
    password,
    sega_id,
    sega_password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
where username = $1;

-- name: GetUserByMaiID :one
select
    user_id,
    username,
    password,
    sega_id,
    sega_password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
where game_name = $1 and tag_line = $2;

-- name: UpdateUser :one
update users
set
    username = $2,
    password = $3,
    sega_id = $4,
    sega_password = $5,
    game_name = $6,
    tag_line = $7,
    updated_at = now()
where user_id = $1
returning user_id, username, password, sega_id, sega_password, game_name, tag_line, updated_at, created_at;
