-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
     id INT PRIMARY KEY,
     short_url VARCHAR(1024) NOT NULL,
     original_url VARCHAR(1024) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS links;
-- +goose StatementEnd