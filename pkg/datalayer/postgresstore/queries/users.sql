-- name: AddUser :exec
INSERT INTO users (username, email, passwd) VALUES ($1, $2, $3);
