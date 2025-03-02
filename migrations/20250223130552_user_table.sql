-- +goose Up
-- +goose StatementBegin
create table if not exists users
(
    id            integer
        primary key autoincrement,
    email         text     not null unique,
    password_hash text     not null,
    created_at    datetime not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists users;
-- +goose StatementEnd
