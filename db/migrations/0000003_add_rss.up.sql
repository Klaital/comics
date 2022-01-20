CREATE TABLE IF NOT EXISTS rss_items (
    item_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    webcomic_id BIGINT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT false,
    guid VARCHAR(64) NOT NULL,
    title VARCHAR(64) NOT NULL,
    link VARCHAR(64) NOT NULL,


    INDEX rss_owner (user_id),
    INDEX rss_singleitem (user_id, guid)
);
