CREATE TABLE score_conversions (
    id             BIGSERIAL PRIMARY KEY,
    section_type   VARCHAR(50) NOT NULL,
    correct_count  INT NOT NULL,
    scaled_score   INT NOT NULL
);

CREATE INDEX idx_score_conversions_section_type ON score_conversions(section_type);
CREATE UNIQUE INDEX idx_score_conversions_type_count ON score_conversions(section_type, correct_count);
