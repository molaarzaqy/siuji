CREATE TABLE options (
    id           BIGSERIAL PRIMARY KEY,
    public_id    UUID NOT NULL DEFAULT gen_random_uuid(),
    question_id  BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    label        VARCHAR(10) NOT NULL,
    option_text  TEXT NOT NULL,
    position     INT NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_options_public_id ON options(public_id);
CREATE INDEX idx_options_question_id ON options(question_id);
