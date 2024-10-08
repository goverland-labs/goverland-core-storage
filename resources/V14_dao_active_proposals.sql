alter table daos
    add active_proposals_ids jsonb default '[]' not null;

update daos
set active_votes = cnt.active_votes, active_proposals_ids = cnt.list
from (select
          dao_id,
          count(id) filter (where state = 'active' and spam is not true) as active_votes,
          coalesce(json_agg(id) filter (where state = 'active' and spam is not true), '[]')  as list
      from proposals
      group by dao_id) cnt
where daos.id = cnt.dao_id;
