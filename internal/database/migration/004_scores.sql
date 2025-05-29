-- +goose Up
create table scores (
    id uuid primary key,
    beatmap_id uuid not null references beatmaps (id) on delete cascade,
    song_id uuid not null references songs (id) on delete cascade,
    user_uuid uuid not null references users (id) on delete cascade,
    accuracy text not null,
    max_combo int not null,
    dx_score int not null,
    tap_critical int not null default 0,
    tap_perfect int not null default 0,
    tap_great int not null default 0,
    tap_good int not null default 0,
    tap_miss int not null default 0,
    hold_critical int not null default 0,
    hold_perfect int not null default 0,
    hold_great int not null default 0,
    hold_good int not null default 0,
    hold_miss int not null default 0,
    slide_critical int not null default 0,
    slide_perfect int not null default 0,
    slide_great int not null default 0,
    slide_good int not null default 0,
    slide_miss int not null default 0,
    touch_critical int not null default 0,
    touch_perfect int not null default 0,
    touch_great int not null default 0,
    touch_good int not null default 0,
    touch_miss int not null default 0,
    break_critical int not null default 0,
    break_perfect int not null default 0,
    break_great int not null default 0,
    break_good int not null default 0,
    break_miss int not null default 0,
    fast int not null default 0,
    late int not null default 0,
    played_at timestamp not null,
    created_at timestamp default now()
);
create index idx_scores_user_uuid on scores (user_uuid);
create index idx_scores_beatmap_id on scores (beatmap_id);
create index idx_scores_song_id on scores (song_id);

-- +goose Down
drop table if exists scores;
drop index if exists idx_scores_user_uuid;
drop index if exists idx_scores_song_id;
drop index if exists idx_scores_beatmap_id;
