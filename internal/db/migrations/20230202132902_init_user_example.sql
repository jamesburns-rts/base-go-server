-- +goose Up
create table example.user_example
(
    user_id      bigint not null references example.app_user (id),
    example_name text   not null,
    primary key (user_id, example_name)
);

-- +goose Down
drop table example.user_example;
