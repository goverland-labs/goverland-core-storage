delete from ens_names en1 using ens_names en2 where en1.name = en2.name and en1.created_at < en2.created_at;

create unique index if not exists ens_names_unique_idx on ens_names(name);
