ALTER TABLE daos
    ADD COLUMN IF NOT EXISTS active_votes integer default 0,
    ADD COLUMN IF NOT EXISTS verified boolean;
