ALTER TABLE proposals
    ADD COLUMN original_state text;

UPDATE proposals
SET original_state = state
WHERE 1 = 1;
