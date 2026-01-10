-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA gophermart;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA gophermart CASCADE;
-- +goose StatementEnd
