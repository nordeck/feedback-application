-- +goose Up
create table tokens
(
    id             serial primary key,
    oidcToken      varchar(1024),
    created_at     timestamp  not null,
);

-- +goose Down
drop table feedbacks;