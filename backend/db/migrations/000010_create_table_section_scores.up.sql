CREATE TABLE section_scores (
    id                     BIGSERIAL PRIMARY KEY,
    public_id              UUID NOT NULL DEFAULT gen_random_uuid(),
    participant_period_id  BIGINT NOT NULL REFERENCES participant_periods(id) ON DELETE CASCADE,
    section_id             BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    correct_count          INT NOT NULL DEFAULT 0,
    raw_score              INT NOT NULL DEFAULT 0,
    scaled_score           INT NOT NULL DEFAULT 0,
    created_at             TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_section_scores_public_id ON section_scores(public_id);
CREATE INDEX idx_section_scores_participant_period_id ON section_scores(participant_period_id);
CREATE INDEX idx_section_scores_section_id ON section_scores(section_id);
CREATE UNIQUE INDEX idx_section_scores_period_section ON section_scores(participant_period_id, section_id);
