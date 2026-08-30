CREATE TABLE participant_answers (
    id                     BIGSERIAL PRIMARY KEY,
    public_id              UUID NOT NULL DEFAULT gen_random_uuid(),
    participant_period_id  BIGINT NOT NULL REFERENCES participant_periods(id) ON DELETE CASCADE,
    question_id            BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_id              BIGINT NOT NULL REFERENCES options(id) ON DELETE CASCADE,
    is_correct             BOOLEAN NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_participant_answers_public_id ON participant_answers(public_id);
CREATE INDEX idx_participant_answers_participant_period_id ON participant_answers(participant_period_id);
CREATE INDEX idx_participant_answers_question_id ON participant_answers(question_id);
CREATE INDEX idx_participant_answers_option_id ON participant_answers(option_id);
CREATE UNIQUE INDEX idx_participant_answers_period_question ON participant_answers(participant_period_id, question_id);
