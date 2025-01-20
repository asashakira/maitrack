-- name: CreateUserScrapeMetadata :one
insert into user_scrape_metadata (
    user_id,
    last_played_at,
    updated_at,
    created_at
)
values ($1, $2, now(), now())
returning user_id, last_played_at, updated_at, created_at;

-- name: GetUserScrapeMetadataByUserID :one
select
    user_id,
    last_played_at,
    updated_at,
    created_at
from user_scrape_metadata
where user_id = $1;

-- name: UpdateUserScrapeMetadata :one
update user_scrape_metadata
set
    last_played_at = $2,
    updated_at = now()
where user_id = $1
returning user_id, last_played_at, updated_at, created_at;
