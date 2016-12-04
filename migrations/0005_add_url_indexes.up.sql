CREATE TABLE saves_url_index (
    url     text,
    save_id uuid,
    UNIQUE (url, save_id)
);
