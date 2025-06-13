-- name: GetConfigs :many
SELECT key, value FROM kvstore;

-- name: FindMatchedWord :one
SELECT exp FROM dictionary WHERE word = ? AND translator = ? AND lang = ?;

-- name: SaveWord :exec
INSERT INTO dictionary (word, exp, translator, lang, context) VALUES (?, ?, ?, ?, ?);

-- name: ChangeConfig :exec
INSERT OR REPLACE INTO  kvstore (key, value) VALUES (?, ?);

-- name: ListStoredLanguages :many
SELECT DISTINCT lang FROM dictionary;
