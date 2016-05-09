CREATE TABLE consumers (
    api_key         uuid        PRIMARY KEY,
    public_key      text        NOT NULL,
    owner_user_id   uuid        ,
    created_at      timestamp   ,
    updated_at      timestamp   ,
    is_internal     bool        ,
    active          bool
);

CREATE INDEX idx_consumers_by_user ON consumers (owner_user_id);
