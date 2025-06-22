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

-- name: ListStoredHistoryContexts :many
SELECT DISTINCT context FROM dictionary;

-- name: ListStoredWords :many
SELECT word, exp FROM dictionary WHERE word LIKE CAST(sqlc.arg(word_like) AS TEXT) || '%' COLLATE NOCASE AND
    (translator = sqlc.arg(translator) OR sqlc.arg(translator) = '') AND
    (lang = sqlc.arg(lang) OR sqlc.arg(lang) = '') AND
    (context = sqlc.arg(context) OR sqlc.arg(context) = '');

-- name: ListStoredNoteContexts :many
SELECT DISTINCT context FROM notes;

-- name: SaveNote :exec
INSERT INTO notes (note, comment, occurrence, context)
VALUES (?1, ?2, 1, ?3)
ON CONFLICT (note) DO UPDATE SET comment = ?2, context = ?3, occurrence = occurrence + 1;

-- name: ListStoredNotes :many
SELECT * FROM notes WHERE note LIKE '%' COLLATE NOCASE || CAST(?1 AS TEXT) || '%' COLLATE NOCASE
AND (context = ?2 OR ?2 = '')
ORDER BY occurrence DESC;
