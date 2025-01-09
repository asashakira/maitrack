-- +goose Up
create table scores (
    score_id uuid primary key,
    beatmap_id uuid not null,
    song_id uuid not null,
    user_id uuid not null,
    accuracy text not null,
    max_combo int not null,
    dx_score int not null,
    played_at timestamp not null,
    created_at timestamp not null
);

-- +goose Down
drop table if exists scores;
