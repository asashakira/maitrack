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
    is_utage bool not null default false,
    is_available bool not null default true,
    is_new bool not null default false,
    release_date date,
    delete_date date,
    updated_at timestamp default now(),
    created_at timestamp default now(),
    unique (title, artist)
);

-- +goose Down
drop table if exists songs;
