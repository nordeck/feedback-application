-- +goose Up
create table feedbacks
(
    id             serial primary key,
    rating         numeric(1) not null,
    rating_comment varchar(1024),
    created_at     timestamp  not null,
    metadata       jsonb,
    jwt            varchar(512)
);

-- +goose Down
drop table feedbacks;