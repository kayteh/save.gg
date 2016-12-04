CREATE TABLE file_entities (
    file_id         uuid,
    url             text,
    alt_url         text,
    owner_id        uuid,

    created_at      timestamp,
    updated_at      timestamp,

    accepted        bool DEFAULT false,
    accepted_at     timestamp,
    rejection_log   text
);

