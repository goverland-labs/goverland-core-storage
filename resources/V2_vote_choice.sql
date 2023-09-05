drop index idx_votes_proposal_id;
drop table votes;

create table votes
(
    id             text not null primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    dao_id         uuid,
    proposal_id    text,
    ipfs           text,
    voter          text,
    created        bigint,
    reason         text,
    choice         jsonb,
    app            text,
    vp             float,
    vp_by_strategy jsonb,
    vp_state       text
);

create index idx_votes_proposal_id on votes (proposal_id);