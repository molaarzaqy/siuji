CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id  BIGSERIAL PRIMARY KEY,
    public_id  UUID NOT NULL DEFAULT gen_random_uuid(),
    email      VARCHAR(255) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(50) NOT NULL DEFAULT 'participant',
    nim        VARCHAR(100),
    university VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_users_public_id ON users(public_id);
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_nim ON users(nim) WHERE nim IS NOT NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);