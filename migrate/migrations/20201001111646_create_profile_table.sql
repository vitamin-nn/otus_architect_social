-- +goose Up
CREATE TABLE user_profile (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255) NOT NULL,
    birthdate date NOT NULL,
    sex enum('M', 'F'),
    interest_list varchar(1024),
    city varchar(255)
);

-- +goose Down
DROP TABLE user_profile;
