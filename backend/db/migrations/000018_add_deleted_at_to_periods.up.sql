ALTER TABLE periods
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_periods_deleted_at ON periods(deleted_at);
