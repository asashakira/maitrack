-- +goose Up
create table users (
    user_id uuid primary key,
    sega_id text not null unique,
    password text not null,
    game_name text not null,
    tag_line text not null,
    updated_at timestamp not null,
    created_at timestamp not null
);

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

create table user_scrape_metadata (
    user_id uuid primary key references users (user_id) on delete cascade,
    last_played_at timestamp not null,
    updated_at timestamp not null,
    created_at timestamp not null
);

-- +goose Down
drop table if exists users;
drop table if exists user_data cascade;
drop table if exists user_scrape_metadata cascade;
