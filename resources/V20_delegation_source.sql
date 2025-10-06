ALTER TABLE delegates_history
    ADD COLUMN IF NOT EXISTS source text default 'split-delegation';

ALTER TABLE delegates_summary
    ADD COLUMN IF NOT EXISTS type text default 'split-delegation';

ALTER TABLE delegates_summary
    ADD COLUMN IF NOT EXISTS chain_id text default null;

WITH strategies AS (SELECT d.id,
                           d.original_id,
                           s ->> 'Name'    AS strategy_name,
                           s ->> 'Network' AS strategy_network,
                           CASE
                               WHEN s ->> 'Name' = 'split-delegation' THEN 1
                               WHEN s ->> 'Name' = 'delegation' THEN 2
                               ELSE 3
                               END         AS priority
                    FROM daos d
                             CROSS JOIN LATERAL jsonb_array_elements(d.strategies) s
                    WHERE d.id IN (SELECT dao_id::uuid FROM delegates_summary GROUP BY dao_id)),
     ranked AS (SELECT id,
                       original_id,
                       strategy_name,
                       strategy_network,
                       ROW_NUMBER() OVER (PARTITION BY id ORDER BY priority) AS rn
                FROM strategies),
     chosen AS (SELECT id,
                       strategy_name,
                       strategy_network
                FROM ranked
                WHERE rn = 1)
UPDATE delegates_summary ds
SET type     = CASE
                   WHEN c.strategy_name = 'split-delegation' THEN 'split-delegation'
                   WHEN c.strategy_name = 'delegation' THEN 'delegation'
                   ELSE 'unrecognized'
    END,
    chain_id = CASE
                   WHEN c.strategy_name != 'split-delegation' THEN c.strategy_network
        END
FROM chosen c
WHERE ds.dao_id::uuid = c.id;

drop index if exists idx_delegates_summary_unique;

create unique index if not exists idx_delegates_summary_unique
    on delegates_summary (address_from, address_to, dao_id, chain_id);

create index concurrently votes_dao_case_insensitive_voter_idx
    on votes (dao_id, lower(voter));

create table erc20_event_history
(
    id              text
        primary key,
    original_dao_id text not null,
    chain_id        text not null,
    block_number    integer,
    log_index       integer,
    type            text not null check (type in ('delegation', 'vp_changes', 'transfer')),
    created_at      timestamp default now(),
    payload         json,
    constraint idx_unique_erc20_event_history
        unique (chain_id, block_number, log_index)
);

create table erc20_delegates
(
    id              bigserial
        primary key,
    address         text           not null,
    dao_id          uuid           not null,
    chain_id        text           not null,
    vp              NUMERIC(78, 0) NOT NULL default 0,
    block_number    integer,
    log_index       integer,
    represented_cnt integer        not null default 0,
    created_at      timestamp               default now(),
    updated_at      timestamp               default now(),
    constraint idx_unique_erc20_delegates
        unique (address, dao_id, chain_id)
);

CREATE INDEX idx_erc20_delegates_dao_chain
    ON erc20_delegates (dao_id, chain_id);

CREATE INDEX idx_erc20_delegates_address
    ON erc20_delegates (address);

create table erc20_balances
(
    id         bigserial
        primary key,
    address    text           not null,
    dao_id     uuid           not null,
    chain_id   text,
    value      NUMERIC(78, 0) NOT NULL default 0,
    created_at timestamp               default now(),
    updated_at timestamp               default now(),
    constraint idx_unique_erc20_balance
        unique (address, dao_id, chain_id)
);

CREATE INDEX idx_erc20_balances_dao_chain
    ON erc20_balances (dao_id, chain_id);

CREATE INDEX idx_erc20_balances_address
    ON erc20_balances (address);

create table erc20_vp_totals
(
    id         bigserial
        primary key,
    dao_id     uuid           not null,
    chain_id   text           not null,
    vp         NUMERIC(78, 0) NOT NULL default 0,
    created_at timestamp               default now(),
    updated_at timestamp               default now(),
    constraint idx_unique_erc20_vp_totals
        unique (dao_id, chain_id)
);
