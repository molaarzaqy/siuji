CREATE TABLE questions (
    id          BIGSERIAL PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid(),
    section_id  BIGINT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    question    TEXT NOT NULL,
    audio_url   TEXT,
    image_url   TEXT,
    passage     TEXT,
    position    INT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_questions_public_id ON questions(public_id);
CREATE INDEX idx_questions_section_id ON questions(section_id);
