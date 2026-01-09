-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart.balance (
    user_id BIGSERIAL NOT NULL REFERENCES gophermart.users(id) ON DELETE CASCADE,
    updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    current NUMERIC(10, 2) NOT NULL DEFAULT 0.0,
    total_withdrawn NUMERIC(10, 2) NOT NULL DEFAULT 0.0
);

CREATE TABLE gophermart.withdrawals (
    user_id BIGSERIAL NOT NULL REFERENCES gophermart.users(id) ON DELETE CASCADE,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_id TEXT REFERENCES gophermart.orders("number") ON DELETE NO ACTION,
    amount NUMERIC(10, 2) NOT NULL DEFAULT 0.0
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart.balance;

DROP TABLE gophermart.withdrawals;
-- +goose StatementEnd
