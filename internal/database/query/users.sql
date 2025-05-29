-- name: CreateUser :one
insert into users (
    id,
    user_id,
    display_name,
    password_hash,
    encrypted_sega_id,
    encrypted_sega_password,
    last_played_at,
    last_scraped_at
)
values (
    $1, $2, $3, $4, $5, $6, $7, $8
)
returning *;


-- name: GetAllUsers :many
select
    id,
    user_id,
    display_name,
    encrypted_sega_id,
    encrypted_sega_password,
    last_played_at
from users;


-- name: GetUserByID :one
select
    u.id,
    u.user_id,
    u.display_name,
    u.last_played_at,
    u.last_scraped_at,
    u.scrape_status,

    d.rating,
    d.season_play_count,
    d.total_play_count,

    m.bio,
    m.profile_image_url,
    m.location,
    m.twitter_id

from users u
left join (
    select distinct on (user_uuid) *
    from user_data
    order by user_uuid, created_at desc
) as d on u.id = d.user_uuid
left join user_metadata m on u.id = m.user_uuid
where u.id = $1;


-- name: GetUserByUserID :one
select
    u.id,
    u.user_id,
    u.display_name,
    u.last_played_at,
    u.last_scraped_at,
    u.scrape_status,

    d.rating,
    d.season_play_count,
    d.total_play_count,

    m.bio,
    m.profile_image_url,
    m.location,
    m.twitter_id

from users u
left join (
    select distinct on (user_uuid) *
    from user_data
    order by user_uuid, created_at desc
) as d on u.id = d.user_uuid
left join user_metadata m on u.id = m.user_uuid
where u.user_id = $1;


-- name: GetPasswordHashByUserID :one
select
    user_id,
    display_name,
    password_hash
from users
where user_id = $1;

-- name: GetSegaCredentialsByUserID :one
select
    encrypted_sega_id,
    encrypted_sega_password
from users
where user_id = $1;


-- name: UpdateUserByUUID :one
update users
set
    user_id = $2,
    display_name = $3,
    password_hash = $4,
    encrypted_sega_id = $5,
    encrypted_sega_password = $6,
    updated_at = now()
where id = $1
returning *;


-- name: UpdateLastPlayedAt :one
update users
set
    last_played_at = $2,
    updated_at = now()
where user_id = $1
returning user_id, last_played_at;


-- name: UpdateLastScrapedAt :one
update users
set
    last_scraped_at = $2,
    updated_at = now()
where user_id = $1
returning user_id, last_scraped_at;
