-- name: CreateUser :one
insert into users (
    user_id,
    sega_id,
    password,
    game_name,
    tag_line,
    updated_at,
    created_at
)
values ($1, $2, $3, $4, $5, now(), now())
returning user_id, sega_id, password, game_name, tag_line, updated_at, created_at;

-- name: GetUserByID :one
select
    user_id,
    sega_id,
    password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
where user_id = $1;

-- name: GetUserBySegaID :one
select
    user_id,
    sega_id,
    password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
where sega_id = $1;

-- name: GetUserByMaiID :one
select
    user_id,
    sega_id,
    password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
where game_name = $1 and tag_line = $2;

-- name: GetAllUsers :many
select
    user_id,
    sega_id,
    password,
    game_name,
    tag_line,
    updated_at,
    created_at
from users
order by updated_at desc;

-- name: UpdateUser :one
update users
set
    sega_id = $1,
    password = $2,
    game_name = $3,
    tag_line = $4,
    updated_at = now()
where user_id = $5
returning user_id, sega_id, password, game_name, tag_line, updated_at, created_at;
