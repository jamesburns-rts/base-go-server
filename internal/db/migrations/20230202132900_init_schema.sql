-- +goose Up
create schema example;

-- +goose Down
drop schema example;
