CREATE INDEX CONCURRENTLY IF NOT EXISTS daos_fulltext_search_idx ON daos USING GIN(to_tsvector('english', name));

CREATE INDEX CONCURRENTLY IF NOT EXISTS proposals_fulltext_search_idx ON proposals USING GIN(to_tsvector('english', title));