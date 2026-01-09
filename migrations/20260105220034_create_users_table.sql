-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart.users (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    "login" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart.users;
-- +goose StatementEnd
