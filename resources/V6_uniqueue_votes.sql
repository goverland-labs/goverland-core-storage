begin transaction;

-- create temporary table
create table if not EXISTS tmp_voter_duplicates
(
    proposal_id text,
    voter       text,
    last_id     text,
    cnt         int
);


-- insert duplicates to temporary table
-- since 1st Sep 00:00 = 1693526400
insert into tmp_voter_duplicates (
    proposal_id,
    voter,
    cnt) (
    select
        proposal_id,
        voter,
        count(*) cnt
    from votes
    where created >= 1693526400
    group by proposal_id, voter
    having (count(*)) >= 2
);

-- prefill last insert id to avoid deleting last revision
do
$$
    declare
        rec record;
    begin
        for rec in select tmp_voter_duplicates.proposal_id, tmp_voter_duplicates.voter from tmp_voter_duplicates
            loop
                update tmp_voter_duplicates
                set last_id = (select id
                               from votes
                               where proposal_id = rec.proposal_id
                                 and voter = rec.voter
                               order by created_at desc
                               limit 1)
                where proposal_id = rec.proposal_id
                  and voter = rec.voter;
            end loop;
    end;
$$;

-- delete duplicates except last revision
do
$$
    declare
        rec record;
    begin
        for rec
            in select
                   tmp_voter_duplicates.proposal_id,
                   tmp_voter_duplicates.voter,
                   tmp_voter_duplicates.last_id
               from tmp_voter_duplicates
            loop
                delete from votes
                where proposal_id = rec.proposal_id
                  and voter = rec.voter
                  and id != rec.last_id;
            end loop;
    end;
$$;

drop table tmp_voter_duplicates;

-- create index
alter table votes
    add constraint votes_unique_proposal_voter
        unique (proposal_id, voter);

commit;
