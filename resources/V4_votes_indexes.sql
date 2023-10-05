drop index if exists idx_votes_proposal_id;

create index idx_votes_proposal_id
    on votes (proposal_id asc, created desc);
