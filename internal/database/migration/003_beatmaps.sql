-- +goose Up
create table beatmaps (
    id uuid primary key,
    song_id uuid not null references songs (id) on delete cascade,
    difficulty text not null, -- basic advanced expert master remaster utage
    level text not null,
    internal_level numeric(3, 1),
    type text not null, -- dx or std or utage
    total_notes int not null default 0,
    tap int not null default 0,
    hold int not null default 0,
    slide int not null default 0,
    touch int not null default 0,
    break int not null default 0,
    note_designer text not null default '?',
    max_dx_score int not null default 0,
    updated_at timestamp default now(),
    created_at timestamp default now(),
    unique (id, difficulty, type)
);
create index idx_beatmaps_song_id on beatmaps (song_id);
create index idx_beatmaps_difficulty_type on beatmaps (difficulty, type);

-- +goose Down
drop table if exists beatmaps;
