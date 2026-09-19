ALTER TABLE periods DROP COLUMN IF EXISTS certificate_public_id;
ALTER TABLE questions DROP COLUMN IF EXISTS audio_public_id;
ALTER TABLE questions DROP COLUMN IF EXISTS image_public_id;
ALTER TABLE participant_periods DROP COLUMN IF EXISTS certificate_public_id;