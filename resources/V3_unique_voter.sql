create table dao_voter
(
    dao_id     uuid primary key,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    voter      text
);

alter table dao_voter
    add constraint idx_dao_voter_unique
        unique (dao_id, voter);

alter table daos
    add members_count integer default 0;

insert into dao_voter (dao_id, created_at, updated_at, voter)
    (select dao_id, now(), now(), voter
     from votes
     group by dao_id, voter)
on conflict DO NOTHING;

update daos
set members_count = cnt.members_count
from (select dao_id, count(*) as members_count
      from dao_voter
      group by dao_id) cnt
where daos.id = cnt.dao_id;
