-- +goose Up
ALTER TABLE user_profile ADD UNIQUE INDEX uniq_index_email (email);

-- +goose Down
ALTER TABLE user_profile DROP INDEX uniq_index_email;
