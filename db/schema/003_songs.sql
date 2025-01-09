-- +goose Up
create table songs (
    song_id uuid primary key,
    alt_key text not null,
    title text not null,
    artist text not null,
    genre text not null,
    bpm text not null,
    image_url text not null,
    version text not null,
    is_utage bool not null,
    is_available bool not null,
    release_date date not null,
    delete_date date,
    updated_at timestamp not null,
    created_at timestamp not null,
    unique (title, artist)
);

-- +goose Down
drop table if exists songs;
