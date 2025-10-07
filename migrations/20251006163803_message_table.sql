-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images(
    id UUID PRIMARY KEY,
    format TEXT CHECK (format IN ('jpg', 'jpeg', 'png', 'gif')),
    status TEXT CHECK (status IN ('in progress', 'finished')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
