-- name: CreateUserMetadata :one
insert into user_metadata (
    user_uuid,
    bio,
    profile_image_url,
    location,
    twitter_id
)
values ($1, $2, $3, $4, $5)
returning *;


-- name: GetUserMetadataByUserID :one
select *
from user_metadata
where user_uuid = $1;


-- name: UpdateUserMetadata :one
update user_metadata
set
    bio = $2,
    profile_image_url = $3,
    location = $4,
    twitter_id = $5,
    updated_at = now()
where user_uuid = $1
returning *;

-- name: UpdateProfileImageUrl :exec
update user_metadata
set
    profile_image_url = $2,
    updated_at = now()
where user_uuid = $1;
