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
select *
from users;

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
