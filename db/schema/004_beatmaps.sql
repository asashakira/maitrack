-- +goose Up
create table beatmaps (
    beatmap_id uuid primary key,
    song_id uuid not null,
    difficulty text not null,
    level text not null,
    internal_level numeric(3, 1) not null default -1,
    type text not null,
    total_notes int not null default 0,
    tap int not null default 0,
    hold int not null default 0,
    slide int not null default 0,
    touch int not null default -1,
    break int not null default 0,
    note_designer text not null default '-',
    max_dx_score int not null,
    is_valid bool not null,
    updated_at timestamp not null,
    created_at timestamp not null
);

-- +goose Down
drop table if exists beatmaps;
