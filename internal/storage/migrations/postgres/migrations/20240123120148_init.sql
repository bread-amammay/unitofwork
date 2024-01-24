-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id         UUID PRIMARY KEY UNIQUE  NOT NULL,
    username   text                     NOT NULL,
    first_name text                     NOT NULL,
    last_name  text                     NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

create table blog_posts
(
    id         UUID PRIMARY KEY UNIQUE  NOT NULL DEFAULT uuid_generate_v4(),
    title      text                     NOT NULL,
    body       text                     NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    user_id    UUID                     NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

drop table blog_posts;
drop table users;

