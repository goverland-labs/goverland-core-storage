create table if not exists custom_dao_recommendations (
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
    dao_id uuid not null,
    strategy_name text not null,
    symbol text not null,
    network_id text,
    address text not null
);

create index if not exists custom_dao_recommendations_dao_id_idx on custom_dao_recommendations(dao_id);
