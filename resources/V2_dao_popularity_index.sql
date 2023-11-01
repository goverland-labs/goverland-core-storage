ALTER TABLE daos
    ADD COLUMN popularity_index double precision;

UPDATE daos
SET popularity_index = voters_count/10000.0
WHERE 1 = 1;
