-- name: CreateUserMetadata :one
insert into user_metadata (
    user_id,
    last_played_at,
    last_scraped_at,
    updated_at,
    created_at
)
values ($1, $2, $3, now(), now())
returning user_id, last_played_at, last_scraped_at, updated_at, created_at;

-- name: GetUserMetadataByUserID :one
select
    user_id,
    last_played_at,
    last_scraped_at,
    updated_at,
    created_at
from user_metadata
where user_id = $1;

-- name: UpdateUserMetadata :one
update user_metadata
set
    last_played_at = $2,
    last_scraped_at = $3,
    updated_at = now()
where user_id = $1
returning user_id, last_played_at, last_scraped_at;

-- name: UpdateLastPlayedAt :one
update user_metadata
set
    last_played_at = $2,
    updated_at = now()
where user_id = $1
returning user_id, last_played_at;

-- name: UpdateLastScrapedAt :one
update user_metadata
set
    last_scraped_at = $2,
    updated_at = now()
where user_id = $1
returning user_id, last_scraped_at;
