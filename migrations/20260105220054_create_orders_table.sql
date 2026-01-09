-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart.orders (
    user_id BIGSERIAL REFERENCES gophermart.users(id) ON DELETE CASCADE,
    "number" TEXT NOT NULL PRIMARY KEY,
    "status" TEXT NOT NULL DEFAULT 'NEW',
    accrual NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart.orders;
-- +goose StatementEnd
