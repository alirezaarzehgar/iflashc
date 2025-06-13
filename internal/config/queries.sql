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

-- name: ListStoredContexts :many
SELECT DISTINCT context FROM dictionary;

-- name: ListStoredWords :many
SELECT word, exp FROM dictionary WHERE translator = ? AND lang = ? AND context = ? AND  word LIKE CAST(sqlc.arg(word_like) AS TEXT) || '%' COLLATE NOCASE;
