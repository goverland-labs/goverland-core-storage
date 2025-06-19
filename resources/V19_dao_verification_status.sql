ALTER TABLE daos
    ADD COLUMN IF NOT EXISTS verification_status text;
ALTER TABLE daos
    ADD COLUMN IF NOT EXISTS verification_comment text;
