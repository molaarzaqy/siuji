CREATE TABLE participant_periods (
    id          BIGSERIAL PRIMARY KEY,
    public_id   UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period_id   BIGINT NOT NULL REFERENCES periods(id) ON DELETE CASCADE,
    status      VARCHAR(50) NOT NULL DEFAULT 'registered'
        CHECK (status IN ('registered', 'started', 'completed')),
    score       INT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_participant_periods_public_id ON participant_periods(public_id);
CREATE INDEX idx_participant_periods_user_id ON participant_periods(user_id);
CREATE INDEX idx_participant_periods_period_id ON participant_periods(period_id);
CREATE UNIQUE INDEX idx_participant_periods_user_period ON participant_periods(user_id, period_id);
