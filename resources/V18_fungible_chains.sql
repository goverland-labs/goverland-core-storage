CREATE TABLE fungible_chains
(
    fungible_id TEXT      NOT NULL,
    chain_id    TEXT      NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    external_id TEXT,
    chain_name  TEXT,
    icon_url    TEXT,
    address     TEXT,
    decimals    INTEGER,
    PRIMARY KEY (fungible_id, chain_id)
);
