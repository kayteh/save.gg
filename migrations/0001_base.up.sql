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

CREATE TABLE sessions (
    session_id      uuid        PRIMARY KEY,
    session_key     uuid        NOT NULL,
    created_at      timestamp   NOT NULL
);

CREATE INDEX idx_session_key ON sessions (session_key)
