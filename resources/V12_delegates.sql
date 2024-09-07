create table delegates_history
(
    action            text not null,
    address_from      text not null,
    chain_id          text,
    original_space_id text not null,
    block_number      integer,
    block_timestamp   integer,
    payload           jsonb,
    created_at        timestamp default now()
);

create index delegates_address_from_index
    on delegates_history (lower(address_from));

create table delegates_summary
(
    address_from         text                    not null,
    address_to           text                    not null,
    dao_id               text                    not null,
    weight               integer,
    expires_at           integer,
    last_block_timestamp integer,
    created_at           timestamp default now() not null
);

create index delegates_summary_address_from_index
    on delegates_summary (lower(address_to));

create unique index if not exists idx_delegates_summary_unique
    on delegates_summary(address_from, address_to, dao_id);
