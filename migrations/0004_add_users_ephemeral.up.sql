CREATE TABLE users_known_ips (
    user_id         uuid,
    ip_address      text,
    last_seen       timestamp,
    UNIQUE(user_id, ip_address)
);

CREATE TABLE users_old_secrets (
    user_id         uuid,
    secret          text,
    changed_at      timestamp,
    changed_by_ip   text
);
