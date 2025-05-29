-- +goose Up
create table users (
    id uuid primary key,
    user_id text not null unique,
    email text unique,
    email_verified boolean default false,
    display_name text not null,
    password_hash text not null,
    encrypted_sega_id text not null unique,
    encrypted_sega_password text not null,
    last_played_at timestamp not null,
    last_scraped_at timestamp not null,
    scrape_status text default 'idle',
    deleted_at timestamp,
    updated_at timestamp default now(),
    created_at timestamp default now()
);
create index idx_users_user_id on users (user_id);

create table user_data (
    id uuid primary key,
    user_id uuid not null references users (id) on delete cascade,
    rating int not null,
    season_play_count int not null,
    total_play_count int not null,
    created_at timestamp default now()
);

create table user_metadata (
    user_id uuid primary key references users (id) on delete cascade,
    bio text,
    profile_image_url text,
    location text,
    twitter_id text,
    updated_at timestamp default now(),
    created_at timestamp default now()
);

-- +goose Down
drop index if exists idx_users_user_id;
drop table if exists user_data cascade;
drop table if exists user_metadata cascade;
drop table if exists users;
