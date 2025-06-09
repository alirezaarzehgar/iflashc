-- name: GetConfigs :many
SELECT key, value FROM kvstore;

-- name: FindMatchedWord :one
SELECT exp FROM dictionary WHERE word = ? AND translator = ?;

-- name: SaveWord :exec
INSERT INTO dictionary (word, exp, translator) VALUES (?, ?, ?);

-- name: ChangeConfig :exec
INSERT OR REPLACE INTO  kvstore (key, value) VALUES (?, ?);
