create index idx_dao_popularity_idx
    on daos (popularity_index desc);

create index idx_dao_category
    on daos (categories);