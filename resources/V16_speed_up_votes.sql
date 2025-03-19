create index if not exists votes_deduplicated_case_insensitive_author_idx
    on votes (lower(voter));
