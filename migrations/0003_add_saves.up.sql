CREATE TABLE saves (

    save_id         uuid        PRIMARY KEY,
    canonical_url   text        UNIQUE NOT NULL,
    custom_url      text        UNIQUE,
    created_at      timestamp   NOT NULL,
    updated_at      timestamp   NOT NULL,
    owner_id        uuid        NOT NULL,

    title           text        NOT NULL,
    description     markdown    ,
    metadeta_id     uuid        ,

    file_entity_id  uuid        UNIQUE NOT NULL

);

CREATE TABLE save_file_entities (
    
);
