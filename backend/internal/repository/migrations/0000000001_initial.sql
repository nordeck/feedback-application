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

CREATE INDEX idx_feedbacks_jwt ON feedbacks(jwt);

-- +goose Down
drop table feedbacks;