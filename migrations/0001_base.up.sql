CREATE TYPE markdown AS (
    markdown    text,
    html        text
); 



CREATE TABLE users (
    user_id     uuid        PRIMARY KEY,
    secret      text        NOT NULL,
    slug        text        UNIQUE NOT NULL,
    username    text        NOT NULL,
    email       text        UNIQUE NOT NULL,
    acl         text[],  
    sub_level   text        DEFAULT 'free',
    activated   bool        DEFAULT false NOT NULL,
    created_at  timestamp   NOT NULL,
    updated_at  timestamp   NOT NULL,
    deleted_at  timestamp,
    session_key uuid        UNIQUE
);



CREATE TABLE saves (
    save_id         uuid        PRIMARY KEY,
    owner_user_id   uuid        NOT NULL,
    url_key         text        UNIQUE NOT NULL,
    created_at      timestamp   NOT NULL,
    updated_at      timestamp   NOT NULL,
    deleted_at      timestamp,
    privacy         text        DEFAULT 'private',
    whitelist       uuid[],
    save_entity_id  uuid        NOT NULL,
    game_id         uuid        NOT NULL,
    title           text,
    description     markdown,   
    metadata        jsonb
);

CREATE INDEX idx_save_metadata ON saves USING gin (metadata);



CREATE TABLE save_entities (
    save_entity_id  uuid        PRIMARY KEY,
    pathname        text        UNIQUE NOT NULL,
    created_at      timestamp   NOT NULL
);



CREATE TABLE comments (
    comment_id      uuid        PRIMARY KEY,
    owner_user_id   uuid        NOT NULL,
    save_id         uuid        NOT NULL,
    body            markdown    NOT NULL,
    created_at      timestamp   NOT NULL,
    updated_at      timestamp   NOT NULL,
    deleted_at      timestamp
);



CREATE TABLE sessions (
    session_id      uuid        PRIMARY KEY,
    session_key     uuid        NOT NULL,
    created_at      timestamp   NOT NULL
);

CREATE INDEX idx_session_key ON sessions (session_key)
