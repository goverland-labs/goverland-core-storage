ALTER TABLE votes
    ADD COLUMN IF NOT EXISTS ens_name text;

create index idx_votes_author
    on votes (voter);
