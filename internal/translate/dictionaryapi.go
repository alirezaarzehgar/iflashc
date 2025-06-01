package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
)

type dictionaryapi struct {
}

type dictionaryResponse []struct {
	Meanings []struct {
		Definitions []struct {
			Definition string   `json:"definition"`
			Synonyms   []string `json:"synonyms"`
			Antonyms   []string `json:"antonyms"`
			Example    string   `json:"example,omitempty"`
		} `json:"definitions"`
		Synonyms []string `json:"synonyms"`
		Antonyms []string `json:"antonyms"`
	} `json:"meanings"`
}

func (d dictionaryapi) Translate(word string) (string, error) {
	urlPath, err := url.JoinPath("https://api.dictionaryapi.dev/api/v2/entries/en/", word)
	if err != nil {
		return "", fmt.Errorf("failt to create url by given word: %w", err)
	}

	resp, err := http.Get(urlPath)
	if err != nil {
		return "", fmt.Errorf("failed to GET request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("%s not found", word)
		}
		return "", fmt.Errorf("response is not successful: %d", resp.StatusCode)
	}

	var dictapiRes dictionaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&dictapiRes); err != nil {
		return "", fmt.Errorf("invalid json format: %w", err)
	}

	meanings := dictapiRes[0].Meanings[0]

	t := template.New("trans table")
	t.Parse(`
# Definitions
{{range .Definitions}}
- {{.Definition}}
{{end}}

# Examples
{{range .Definitions}}
{{if .Example}}
- {{ .Example}}
{{end}}
{{end}}

# Synonyms
{{range .Synonyms}}
- {{.}}
{{end}}
`)

	var tpl bytes.Buffer
	err = t.Execute(&tpl, meanings)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.String(), nil
}
