CREATE TABLE answer_keys (
    id           BIGSERIAL PRIMARY KEY,
    public_id    UUID NOT NULL DEFAULT gen_random_uuid(),
    option_id    BIGINT NOT NULL REFERENCES options(id) ON DELETE CASCADE,
    question_id  BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_answer_keys_public_id ON answer_keys(public_id);
CREATE INDEX idx_answer_keys_option_id ON answer_keys(option_id);
CREATE UNIQUE INDEX idx_answer_keys_question_id ON answer_keys(question_id);
