ALTER TABLE sections ADD section_type VARCHAR(20);
ALTER TABLE sections ADD CONSTRAINT sections_section_type_check
    CHECK (section_type IS NULL OR section_type IN ('listening', 'structure', 'reading'));