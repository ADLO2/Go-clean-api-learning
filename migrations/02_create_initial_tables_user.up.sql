CREATE TABLE users
(
    user_id    VARCHAR(250)                     PRIMARY KEY, /* UUID */
    username      VARCHAR(250)             NOT NULL,
    password    VARCHAR(250)             NOT NULL,
    name VARCHAR(250)             NOT NULL,
    created_at TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP                NOT NULL DEFAULT CURRENT_TIMESTAMP
);