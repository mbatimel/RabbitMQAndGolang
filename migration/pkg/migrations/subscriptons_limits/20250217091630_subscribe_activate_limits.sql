-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE SCHEMA IF NOT EXISTS main;
CREATE TABLE IF NOT EXISTS main.subscriptions
(
    id                 int primary key,
    supplier_id        uuid                      not null,
    limit_id           int,                      
    price              int                       not null,
    created_at         timestamptz default now() not null,
);
CREATE TABLE IF NOT EXISTS main.limits
(
    id                 int primary key,
    count              int                       not null,                     
    describe           string,
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
