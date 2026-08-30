CREATE TABLE periods (
    id                      BIGSERIAL PRIMARY KEY,
    public_id               UUID NOT NULL DEFAULT gen_random_uuid(),
    title                   VARCHAR(255) NOT NULL,
    month                   VARCHAR(50),
    year                    INT,
    status                  VARCHAR(50) NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'published', 'closed')),
    certificate_url         TEXT,
    certificate_exp_month   TIMESTAMP,
    min_passing_grade       INT,
    max_passing_grade       INT,
    start_time              TIMESTAMP,
    end_time                TIMESTAMP,
    created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at              TIMESTAMP
);

CREATE UNIQUE INDEX idx_periods_public_id ON periods(public_id);
CREATE INDEX idx_periods_deleted_at ON periods(deleted_at);
CREATE INDEX idx_periods_status ON periods(status);
