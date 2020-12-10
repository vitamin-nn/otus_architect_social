-- +goose Up
ALTER TABLE user_message CHANGE COLUMN sent_at sent_at DATETIME NOT NULL;

-- +goose Down
