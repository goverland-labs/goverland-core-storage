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
                               WHEN s ->> 'Name' = 'erc20-votes' THEN 3
                               ELSE 4
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
SET type     = c.strategy_name,
    chain_id = CASE
                   WHEN c.strategy_name != 'split-delegation' THEN c.strategy_network
        END
FROM chosen c
WHERE ds.dao_id::uuid = c.id;
