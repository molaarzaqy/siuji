ALTER TABLE sections DROP CONSTRAINT IF EXISTS sections_section_type_check;
ALTER TABLE sections DROP COLUMN IF EXISTS section_type;