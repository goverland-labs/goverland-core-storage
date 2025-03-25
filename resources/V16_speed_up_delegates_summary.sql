create index delegates_summary_expires_at_idx
    on delegates_summary (expires_at) where expires_at > 0;