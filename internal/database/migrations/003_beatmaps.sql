-- +goose Up
create table beatmaps (
    beatmap_id uuid primary key,
    song_id uuid not null references songs (song_id) on delete cascade,
    difficulty text not null,
    level text not null,
    internal_level numeric(3, 1),
    type text not null, -- dx or std
    total_notes int not null default 0,
    tap int not null default 0,
    hold int not null default 0,
    slide int not null default 0,
    touch int not null default 0,
    break int not null default 0,
    note_designer text,
    max_dx_score int not null default 0,
    is_valid bool not null default true,
    updated_at timestamp not null,
    created_at timestamp not null,
    unique (song_id, difficulty, type)
);
create index idx_beatmaps_song_id on beatmaps (song_id);
create index idx_beatmaps_difficulty on beatmaps (difficulty);
create index idx_beatmaps_type on beatmaps (type);

-- +goose Down
drop table if exists beatmaps;
