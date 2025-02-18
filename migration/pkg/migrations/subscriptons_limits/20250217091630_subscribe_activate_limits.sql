-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE SCHEMA IF NOT EXISTS main;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS main.subscriptions
(
    id                 uuid primary key default uuid_generate_v4(),
    supplier_id        int                      not null,
    limit_id           int,                      
    price              int                       not null,
    created_at         timestamptz default now() not null
);
CREATE TABLE IF NOT EXISTS main.limits
(
    id                 uuid primary key default uuid_generate_v4(),
    limit_id           int,                      
    count              int                       not null,                     
    describe           text
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
