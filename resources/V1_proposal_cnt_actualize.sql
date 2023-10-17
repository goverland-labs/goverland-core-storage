update daos
set proposals_count = cnt.proposals_cnt
from (select dao_id, count(*) as proposals_cnt
      from proposals
      group by dao_id) cnt
where daos.id = cnt.dao_id;