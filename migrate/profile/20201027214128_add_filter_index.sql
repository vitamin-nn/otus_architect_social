-- +goose Up
CREATE INDEX f_l_name_idx on user_profile(first_name, last_name);
-- +goose Down
DROP INDEX f_l_name_idx;

