CREATE TABLE period_sections (
    id          BIGSERIAL PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid(),
    period_id   BIGINT NOT NULL REFERENCES periods(id) ON DELETE CASCADE,
    section_id  BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    position    INT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_period_sections_public_id ON period_sections(public_id);
CREATE INDEX idx_period_sections_period_id ON period_sections(period_id);
CREATE INDEX idx_period_sections_section_id ON period_sections(section_id);
CREATE UNIQUE INDEX idx_period_sections_period_section ON period_sections(period_id, section_id);
