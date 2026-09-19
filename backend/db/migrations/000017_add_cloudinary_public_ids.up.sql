ALTER TABLE periods ADD COLUMN certificate_public_id VARCHAR(255);
ALTER TABLE questions ADD COLUMN audio_public_id VARCHAR(255);
ALTER TABLE questions ADD COLUMN image_public_id VARCHAR(255);
ALTER TABLE participant_periods ADD COLUMN certificate_public_id VARCHAR(255);