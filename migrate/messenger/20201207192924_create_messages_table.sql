-- +goose Up
CREATE TABLE user_message (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    dialog_id int NOT NULL,
    sender_user_id int NOT NULL,
    receiver_user_id int NOT NULL,
    sent_at date NOT NULL,
    body text NOT NULL,
    is_read tinyint NOT NULL
);

-- +goose Down
DROP TABLE user_message;
