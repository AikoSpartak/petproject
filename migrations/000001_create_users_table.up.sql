CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY,
    full_name     VARCHAR(100) NOT NULL,
    email         VARCHAR(100) NOT NULL UNIQUE,
    phone_number  VARCHAR(30) NOT NULL,
    date_birthday DATE NOT NULL,
    passport      VARCHAR(12) NOT NULL UNIQUE,
    password      VARCHAR(256) NOT NULL,
    role          VARCHAR(15) NOT NULL,
    is_banned     BOOLEAN NOT NULL DEFAULT FALSE,
    country       VARCHAR(256) NOT NULL
);