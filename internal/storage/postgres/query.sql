-- name: SaveURL :one
INSERT INTO urls (URL, shortURL) VALUES ($1, $2) RETURNING shortURL;

-- name: GetURL :one
SELECT URL FROM urls WHERE shortURL = $1;
