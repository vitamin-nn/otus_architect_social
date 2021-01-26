-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_feed (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id int NOT NULL,
    create_at datetime NOT NULL,
    body text NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_feed;
-- +goose StatementEnd
