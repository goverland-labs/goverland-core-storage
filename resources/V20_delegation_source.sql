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
    updated_at      timestamp               default now()
);

CREATE UNIQUE INDEX idx_unique_erc20_delegates
    ON erc20_delegates (lower(address), dao_id, chain_id);

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

create table erc20_totals
(
    id               bigserial
        primary key,
    dao_id           uuid           not null,
    chain_id         text           not null,
    voting_power     NUMERIC(78, 0) NOT NULL default 0,
    total_delegators integer        NOT NULL default 0,
    created_at       timestamp               default now(),
    updated_at       timestamp               default now(),
    constraint idx_unique_erc20_totals
        unique (dao_id, chain_id)
);

ALTER TABLE delegates_summary
    ADD COLUMN IF NOT EXISTS log_index integer default 0;

CREATE INDEX idx_delegates_summary_to_dao_chain_from
    ON delegates_summary (address_to, dao_id, chain_id, address_from);

CREATE INDEX if not exists idx_erc20_event_history_original_dao_id
    ON erc20_event_history (original_dao_id);

CREATE TABLE IF NOT EXISTS erc20_delegations
(
    token                text                    not null,
    chain_id             text                    not null,
    address_from         text                    not null,
    address_to           text                    not null,
    created_at           timestamp default now() not null,
    last_block_timestamp integer,
    log_index            integer   default 0
);

create index idx_erc20_delegations_address_to
    on erc20_delegations (address_to);

create index idx_erc20_delegations_address_from
    on erc20_delegations (address_from);

create unique index if not exists idx_erc20_delegations_unique
    on erc20_delegations (address_from, token, chain_id);

create index if not exists idx_erc20_delegations_to_delegation
    on erc20_delegations (address_to, token, chain_id);

create index if not exists idx_erc20_delegations_token_delegations
    on erc20_delegations (token, chain_id);

truncate erc20_delegations;

insert into erc20_delegations(token, chain_id, address_from, address_to, last_block_timestamp, log_index)
select case
           -- "ens.eth"
           when dao_id = '155f2735-f054-4f59-a0f2-06f05dff1567' then lower('0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72')
           -- "parason.eth"
           when dao_id = 'c5897bb6-9ba7-47dd-9e23-62b60a2ce6e2' then lower('0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91')
           -- "arbitrumfoundation.eth"
           when dao_id = 'c8c85359-9056-4b42-9311-f13de2051389' then lower('0x912ce59144191c1204e64559fe8253a0e49e6548')
           -- "integration-test.eth",
           when dao_id = 'd79cb393-86f4-434f-a26b-7047c11fad11' then lower('0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766')
           end             token,
       chain_id,
       lower(address_from) address_from,
       lower(address_to)   address_to,
       last_block_timestamp,
       log_index
from delegates_summary
where type = 'erc20-votes'
  and dao_id IN (
                 '155f2735-f054-4f59-a0f2-06f05dff1567',
                 'c5897bb6-9ba7-47dd-9e23-62b60a2ce6e2',
                 'd79cb393-86f4-434f-a26b-7047c11fad11',
                 'c8c85359-9056-4b42-9311-f13de2051389');

-- erc20_balances
alter table erc20_balances
    add column token text;

update erc20_balances
set token = case
    -- "ens.eth"
                when dao_id = '155f2735-f054-4f59-a0f2-06f05dff1567'
                    then lower('0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72')

    -- "parason.eth"
                when dao_id = 'c5897bb6-9ba7-47dd-9e23-62b60a2ce6e2'
                    then lower('0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91')

    -- "arbitrumfoundation.eth"
                when dao_id = 'c8c85359-9056-4b42-9311-f13de2051389'
                    then lower('0x912ce59144191c1204e64559fe8253a0e49e6548')
    -- "integration-test.eth",
                when dao_id = 'd79cb393-86f4-434f-a26b-7047c11fad11'
                    then lower('0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766')

                else token
    end
where token is null;
create index idx_erc20_balances_token_chain
    on erc20_balances (token, chain_id);


-- erc20_event_history
alter table erc20_event_history
    add column token text;

update erc20_event_history
set token = case
    -- original_dao_id: "ens.eth"
                when original_dao_id = 'ens.eth'
                    then lower('0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72')
    -- original_dao_id: "parason.eth"
                when original_dao_id = 'parason.eth'
                    then lower('0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91')
    -- original_dao_id: "arbitrumfoundation.eth"
                when original_dao_id = 'parason.eth'
                    then lower('0x912ce59144191c1204e64559fe8253a0e49e6548')
    -- "integration-test.eth",
                when original_dao_id = 'integration-test.eth'
                    then lower('0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766')
                else token
    end
where token is null;

create index idx_erc20_event_history_token
    on erc20_event_history (token);


-- erc20_delegates
alter table erc20_delegates
    add column token text;
update erc20_delegates
set token = case
    -- "ens.eth"
                when dao_id = '155f2735-f054-4f59-a0f2-06f05dff1567'
                    then lower('0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72')

    -- "parason.eth"
                when dao_id = 'c5897bb6-9ba7-47dd-9e23-62b60a2ce6e2'
                    then lower('0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91')

    -- "arbitrumfoundation.eth"
                when dao_id = 'c8c85359-9056-4b42-9311-f13de2051389'
                    then lower('0x912ce59144191c1204e64559fe8253a0e49e6548')
    -- "integration-test.eth",
                when dao_id = 'd79cb393-86f4-434f-a26b-7047c11fad11'
                    then lower('0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766')

                else token
    end
where token is null;

create index idx_erc20_delegates_token_chain
    on erc20_delegates (token, chain_id);


-- erc20_totals
alter table erc20_totals
    add column token text;

update erc20_totals
set token = case
    -- "ens.eth"
                when dao_id = '155f2735-f054-4f59-a0f2-06f05dff1567'
                    then lower('0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72')

    -- "parason.eth"
                when dao_id = 'c5897bb6-9ba7-47dd-9e23-62b60a2ce6e2'
                    then lower('0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91')

    -- "arbitrumfoundation.eth"
                when dao_id = 'c8c85359-9056-4b42-9311-f13de2051389'
                    then lower('0x912ce59144191c1204e64559fe8253a0e49e6548')
    -- "integration-test.eth",
                when dao_id = 'd79cb393-86f4-434f-a26b-7047c11fad11'
                    then lower('0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766')

                else token
    end
where token is null;

create table erc20_token_mapping
(
    dao_id     uuid primary key,
    token      text not null,
    chain_id   text not null,
    created_at timestamp default now()
);

create unique index idx_dao_token_map_token_chain
    on erc20_token_mapping (token, chain_id);


insert into erc20_token_mapping(dao_id, token, chain_id)
values
-- ens.eth
('155f2735-f054-4f59-a0f2-06f05dff1567', lower('0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72'), '1'),
-- arbitrumfoundation.eth
('c8c85359-9056-4b42-9311-f13de2051389', lower('0x912ce59144191c1204e64559fe8253a0e49e6548'), '42161'),
-- parason.eth
('c5897bb6-9ba7-47dd-9e23-62b60a2ce6e2', lower('0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91'), '1');


create materialized view mixed_delegations as

-- split-delegations
select lower(d.address_from) as address_from,
       lower(d.address_to)   as address_to,
       d.dao_id::text        as dao_id,
       d.weight              as weight,
       d.expires_at          as expires_at,
       d.type                as type,
       d.chain_id            as chain_id,
       d.voting_power        as voting_power
from delegates_summary d
where d.type != 'erc20-votes'

union all

-- ERC20 votes
select lower(d.address_from) as address_from,
       lower(d.address_to)   as address_to,
       m.dao_id::text        as dao_id,
       10000::integer        as weight,
       null                  as expires_at,
       'erc20-votes'::text   as type,
       d.chain_id,
       coalesce(ed.vp, 0)    as voting_power
from erc20_delegations d
         join erc20_token_mapping m
              on m.token = d.token
                  and m.chain_id = d.chain_id
         left join erc20_delegates ed
                   on lower(ed.address) = lower(d.address_from)
                       and ed.token = m.token
                       and ed.chain_id = d.chain_id;

create unique index idx_mixed_delegations_unique
    on mixed_delegations (
                          address_from, address_to, dao_id, chain_id, type
        );

create index mixed_delegations_address_to_index
    on mixed_delegations (lower(address_to));

create index idx_mixed_delegations_to_dao_chain_from
    on mixed_delegations (
                          address_to,
                          dao_id,
                          chain_id,
                          address_from
        );

refresh materialized view concurrently mixed_delegations;

alter table erc20_balances
    drop constraint if exists idx_unique_erc20_balance;
alter table erc20_balances
    add constraint idx_unique_erc20_balance_token
        unique (address, token, chain_id);
drop index if exists idx_erc20_balances_dao_chain;

drop index if exists idx_erc20_event_history_original_dao_id;

drop index if exists idx_unique_erc20_delegates;
create unique index idx_unique_erc20_delegates_token
    on erc20_delegates (lower(address), token, chain_id);
drop index if exists idx_erc20_delegates_dao_chain;

alter table erc20_totals
    drop constraint if exists  idx_unique_erc20_totals;
alter table erc20_totals
    add constraint idx_unique_erc20_totals_token
        unique (token, chain_id);

alter table erc20_balances
    drop column dao_id;

alter table erc20_delegates
    drop column dao_id;

alter table erc20_totals
    drop column dao_id;

alter table erc20_event_history
    drop column original_dao_id;

delete
from delegates_summary
where type = 'erc20-votes';
