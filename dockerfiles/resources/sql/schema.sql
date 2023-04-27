create table daos
(
    id              text not null
        primary key,
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
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
    parent_id       text
);

create table proposals
(
    id             text not null
        primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    ipfs           text,
    author         text,
    created        bigint,
    dao_id         text,
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

create table registered_events
(
    id         bigserial
        primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    type       text,
    event      text,
    type_id    text
);

create index idx_registered_events_deleted_at
    on registered_events (deleted_at);

alter table registered_events
    add constraint idx_unique_registered_events
        unique (type, type_id, event);
