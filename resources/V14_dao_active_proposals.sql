alter table daos
    add active_proposals_ids jsonb default '[]' not null;

update daos
set active_votes = cnt.active_votes, active_proposals_ids = cnt.list
from (select
          dao_id,
          count(id) as active_votes,
          json_agg(id) list
      from proposals
      where state = 'active' and spam is not true
      group by dao_id
      having count(id) > 0) cnt
where daos.id = cnt.dao_id;