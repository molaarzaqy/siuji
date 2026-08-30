CREATE TABLE sections (
    id          BIGSERIAL PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid(),
    title       VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_sections_public_id ON sections(public_id);
