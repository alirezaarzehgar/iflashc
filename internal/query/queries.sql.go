// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: queries.sql

package query

import (
	"context"
)

const changeConfig = `-- name: ChangeConfig :exec
INSERT OR REPLACE INTO  kvstore (key, value) VALUES (?, ?)
`

type ChangeConfigParams struct {
	Key   string
	Value string
}

func (q *Queries) ChangeConfig(ctx context.Context, arg ChangeConfigParams) error {
	_, err := q.db.ExecContext(ctx, changeConfig, arg.Key, arg.Value)
	return err
}

const findMatchedWord = `-- name: FindMatchedWord :one
SELECT exp FROM dictionary WHERE word = ? AND translator = ? AND lang = ?
`

type FindMatchedWordParams struct {
	Word       string
	Translator string
	Lang       string
}

func (q *Queries) FindMatchedWord(ctx context.Context, arg FindMatchedWordParams) (string, error) {
	row := q.db.QueryRowContext(ctx, findMatchedWord, arg.Word, arg.Translator, arg.Lang)
	var exp string
	err := row.Scan(&exp)
	return exp, err
}

const getConfigs = `-- name: GetConfigs :many
SELECT key, value FROM kvstore
`

func (q *Queries) GetConfigs(ctx context.Context) ([]Kvstore, error) {
	rows, err := q.db.QueryContext(ctx, getConfigs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Kvstore
	for rows.Next() {
		var i Kvstore
		if err := rows.Scan(&i.Key, &i.Value); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStoredContexts = `-- name: ListStoredContexts :many
SELECT DISTINCT context FROM dictionary
`

func (q *Queries) ListStoredContexts(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listStoredContexts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var context string
		if err := rows.Scan(&context); err != nil {
			return nil, err
		}
		items = append(items, context)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStoredLanguages = `-- name: ListStoredLanguages :many
SELECT DISTINCT lang FROM dictionary
`

func (q *Queries) ListStoredLanguages(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listStoredLanguages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var lang string
		if err := rows.Scan(&lang); err != nil {
			return nil, err
		}
		items = append(items, lang)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStoredWords = `-- name: ListStoredWords :many
SELECT word, exp FROM dictionary WHERE word LIKE CAST(?1 AS TEXT) || '%' COLLATE NOCASE AND
    (translator = ?2 OR ?2 = '') AND
    (lang = ?3 OR ?3 = '') AND
    (context = ?4 OR ?4 = '')
`

type ListStoredWordsParams struct {
	WordLike   string
	Translator string
	Lang       string
	Context    string
}

type ListStoredWordsRow struct {
	Word string
	Exp  string
}

func (q *Queries) ListStoredWords(ctx context.Context, arg ListStoredWordsParams) ([]ListStoredWordsRow, error) {
	rows, err := q.db.QueryContext(ctx, listStoredWords,
		arg.WordLike,
		arg.Translator,
		arg.Lang,
		arg.Context,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListStoredWordsRow
	for rows.Next() {
		var i ListStoredWordsRow
		if err := rows.Scan(&i.Word, &i.Exp); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveWord = `-- name: SaveWord :exec
INSERT INTO dictionary (word, exp, translator, lang, context) VALUES (?, ?, ?, ?, ?)
`

type SaveWordParams struct {
	Word       string
	Exp        string
	Translator string
	Lang       string
	Context    string
}

func (q *Queries) SaveWord(ctx context.Context, arg SaveWordParams) error {
	_, err := q.db.ExecContext(ctx, saveWord,
		arg.Word,
		arg.Exp,
		arg.Translator,
		arg.Lang,
		arg.Context,
	)
	return err
}
