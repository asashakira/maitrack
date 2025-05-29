-- +goose Up
create table songs (
    id uuid primary key,
    alt_key text not null unique,
    title text not null,
    artist text not null,
    genre text not null,
    bpm text not null,
    image_url text not null,
    version text not null,
    sort text not null,
    is_utage bool default false,
    is_available bool default true,
    is_new bool default false,
    release_date date,
    delete_date date,
    updated_at timestamp not null,
    created_at timestamp not null,
    unique (title, artist)
);

-- +goose Down
drop table if exists songs;
