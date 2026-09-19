DROP INDEX IF EXISTS idx_periods_deleted_at;

ALTER TABLE periods
    DROP COLUMN IF EXISTS deleted_at;
