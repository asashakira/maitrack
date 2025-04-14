-- +goose Up
create table users (
    user_id uuid primary key,
    username text not null unique,
    password text not null,
    sega_id text not null unique,
    sega_password text not null,
    game_name text not null,
    tag_line text not null,
    updated_at timestamp not null,
    created_at timestamp not null,
    unique (game_name, tag_line)
);
create index idx_users_username on users (username);
create index idx_users_sega_id on users (sega_id);
create index idx_users_mai_id on users (game_name, tag_line);

create table user_data (
    id uuid primary key,
    user_id uuid not null references users (user_id) on delete cascade,
    game_name text not null,
    tag_line text not null,
    rating int not null,
    season_play_count int not null,
    total_play_count int not null,
    created_at timestamp not null
);

create table user_metadata (
    user_id uuid primary key references users (user_id) on delete cascade,
    last_played_at timestamp not null,
    last_scraped_at timestamp not null,
    updated_at timestamp not null,
    created_at timestamp not null
);

-- +goose Down
drop index if exists idx_users_username;
drop index if exists idx_users_sega_id;
drop index if exists idx_users_mai_id;
drop table if exists user_data cascade;
drop table if exists user_metadata cascade;
drop table if exists users;
