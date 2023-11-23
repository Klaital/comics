CREATE TABLE webcomic (
    webcomic_id SERIAL PRIMARY KEY,
    ordinal INTEGER,
    user_id INTEGER NOT NULL,
    title VARCHAR(128) NOT NULL,
    base_url VARCHAR(128),
    first_comic_url VARCHAR(128),
    latest_comic_url VARCHAR(128),
    rss_url VARCHAR(128),
    updates_monday BOOLEAN NOT NULL DEFAULT 'false',
    updates_tuesday BOOLEAN NOT NULL DEFAULT 'false',
    updates_wednesday BOOLEAN NOT NULL DEFAULT 'false',
    updates_thursday BOOLEAN NOT NULL DEFAULT 'false',
    updates_friday BOOLEAN NOT NULL DEFAULT 'false',
    updates_saturday BOOLEAN NOT NULL DEFAULT 'false',
    updates_sunday BOOLEAN NOT NULL DEFAULT 'false',
    active BOOLEAN NOT NULL default 'true',
    nsfw BOOLEAN NOT NULL default 'false',
    last_read TIMESTAMP
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(32) NOT NULL,
    email VARCHAR(128) NOT NULL UNIQUE,
    passwd VARCHAR(128) NOT NULL
);

CREATE TABLE rss_items (
    item_id SERIAL PRIMARY KEY,
    webcomic_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT 'false',
    guid VARCHAR(128) NOT NULL,
    title VARCHAR(128) NOT NULL,
    link VARCHAR(128) NOT NULL
);
