-- name: AddComic :exec
INSERT INTO webcomic (ordinal, user_id, title,
                      base_url, first_comic_url, latest_comic_url, rss_url,
                      updates_monday, updates_tuesday, updates_wednesday,
                      updates_thursday, updates_friday, updates_saturday, updates_sunday,
                      active, nsfw, last_read)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17);

-- name: AddRssItem :exec
INSERT INTO rss_items (webcomic_id, user_id, guid, title, link)
VALUES ($1, $2, $3, $4, $5);

-- name: ListComics :many
SELECT webcomic_id, ordinal, user_id, title,
       base_url, first_comic_url, latest_comic_url, rss_url,
       updates_monday, updates_tuesday, updates_wednesday,
       updates_thursday, updates_friday, updates_saturday, updates_sunday,
       active, nsfw, last_read
FROM webcomic WHERE user_id = $1;

-- name: MarkItemRead :exec
UPDATE rss_items
SET is_read = $3
WHERE
    user_id = $1
  AND guid = $2;
