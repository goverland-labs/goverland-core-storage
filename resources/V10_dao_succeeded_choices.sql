create table if not exists dao_succeeded_choices(
    dao_id uuid not null primary key,
    choices text[]
);
