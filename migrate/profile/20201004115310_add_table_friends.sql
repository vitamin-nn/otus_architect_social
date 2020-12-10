-- +goose Up
CREATE TABLE user_friends (
    user_id1 int NOT NULL,
    user_id2 int NOT NULL,
    primary key (user_id1, user_id2)
);

-- +goose Down
DROP TABLE user_friends;
