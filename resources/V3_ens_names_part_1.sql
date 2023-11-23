ALTER TABLE proposals
    ADD COLUMN IF NOT EXISTS ens_name text;

CREATE TABLE ens_names
(
    address    TEXT NOT NULL
        PRIMARY KEY,
    name       TEXT,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp
);
