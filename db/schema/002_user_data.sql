-- +goose Up
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

-- +goose Down
drop table if exists user_data;
