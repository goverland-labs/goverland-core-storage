create table daos
(
    id              uuid not null primary key,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    original_id     text,
    name            text,
    private         boolean,
    about           text,
    avatar          text,
    terms           text,
    location        text,
    website         text,
    twitter         text,
    github          text,
    coingecko       text,
    email           text,
    network         text,
    symbol          text,
    skin            text,
    domain          text,
    strategies      jsonb,
    voting          jsonb,
    categories      jsonb,
    treasures       jsonb,
    followers_count integer,
    proposals_count integer,
    guidelines      text,
    template        text,
    parent_id       uuid,
    activity_since  bigint
);

alter table daos
    add constraint daos_idx_unique_original_id unique (original_id);

CREATE INDEX idx_gin_dao_categories ON daos USING gin (categories jsonb_path_ops);
CREATE INDEX idx_dao_name ON daos (lower(name) varchar_pattern_ops);

create table proposals
(
    id             text not null primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    ipfs           text,
    author         text,
    created        bigint,
    dao_id         uuid,
    network        text,
    symbol         text,
    type           text,
    strategies     jsonb,
    title          text,
    body           text,
    discussion     text,
    choices        jsonb,
    start          bigint,
    "end"          bigint,
    quorum         numeric,
    privacy        text,
    snapshot       text,
    state          text,
    link           text,
    app            text,
    scores         jsonb,
    scores_state   text,
    scores_total   numeric,
    scores_updated bigint,
    votes          bigint
);

create index idx_proposal_dao_id on proposals (dao_id);

create table registered_events
(
    id         bigserial primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type       text,
    event      text,
    type_id    text
);

create index idx_registered_events_deleted_at on registered_events (deleted_at);

alter table registered_events
    add constraint idx_unique_registered_events unique (type, type_id, event);

create table votes
(
    id          text not null primary key,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    proposal_id text,
    ipfs        text,
    voter       text,
    created     bigint,
    reason      text
);

create table dao_ids
(
    original_id text not null primary key,
    internal_id uuid not null
);

create index idx_votes_proposal_id on votes (proposal_id);
