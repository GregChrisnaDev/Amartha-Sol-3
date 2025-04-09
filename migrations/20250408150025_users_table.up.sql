CREATE TABLE users
(
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name            VARCHAR(255),
    address         VARCHAR(255),
    password_hash   VARCHAR(255),
    role            INTEGER,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX unq_users_name ON users(name);