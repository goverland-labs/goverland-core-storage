create index votes_proposal_vp_idx
    on votes (proposal_id, vp);

create index votes_proposal_ens_idx
    on votes (proposal_id, ens_name);