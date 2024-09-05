insert into dao_succeeded_choices (dao_id, choices)
    (select internal_id, '{"accept", "enable token transferability", "accept the proposal"}' from dao_ids where original_id = 'safe.eth');