// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package query

type History struct {
	Word       string
	Exp        string
	Lang       string
	Translator string
	Context    string
}

type Kvstore struct {
	Key   string
	Value string
}

type Note struct {
	Note       string
	Comment    string
	Occurrence int64
	Context    string
}
