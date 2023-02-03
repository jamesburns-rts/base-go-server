-- +goose Up
create table example.app_user
(
    id            serial primary key,
    email         text      not null,
    password_hash text      not null,
    first_name    text,
    last_name     text,
    created_at    timestamp not null default now(),
    updated_at    timestamp not null default now(),
    deleted_at    timestamp
);

create unique index on example.app_user (email, (deleted_at is null)) where deleted_at is null;

-- +goose Down
drop table example.app_user;
