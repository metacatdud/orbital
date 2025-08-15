ALTER TABLE applications ADD COLUMN created_at DATETIME;
UPDATE applications SET created_at = CURRENT_TIMESTAMP WHERE created_at IS NULL;

ALTER TABLE applications ADD COLUMN updated_at DATETIME;
UPDATE applications SET updated_at = CURRENT_TIMESTAMP WHERE updated_at IS NULL;

ALTER TABLE applications ADD COLUMN deleted_at DATETIME;
ALTER TABLE applications ADD COLUMN is_enabled BOOLEAN NOT NULL DEFAULT 1;