CREATE TABLE otps (
    id          BIGSERIAL PRIMARY KEY,
    email       VARCHAR(255) NOT NULL,
    code        VARCHAR(10) NOT NULL,
    purpose     VARCHAR(50) NOT NULL,
    expires_at  TIMESTAMP NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP
);

CREATE INDEX idx_otps_email ON otps(email);
CREATE INDEX idx_otps_deleted_at ON otps(deleted_at);
