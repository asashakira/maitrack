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

-- +goose Down
drop table if exists users;
