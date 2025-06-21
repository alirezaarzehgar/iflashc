package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

type fastdic struct {
}

type fastdicSuggestionResponse struct {
	Word    string `json:"word"`
	Meaning string `json:"meaning"`
}

func (f fastdic) Translate(text string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, "https://fastdic.com/suggestions", bytes.NewBuffer([]byte(`{"word": "`+text+`"}`)))
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://fastdic.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://fastdic.com")

	res, err := (&http.Client{Timeout: time.Second * 5}).Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send http request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request is not suggess: %d", res.StatusCode)
	}

	var fsr, meanings []fastdicSuggestionResponse
	err = json.NewDecoder(res.Body).Decode(&fsr)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	for _, r := range fsr {
		if !strings.Contains(r.Meaning, "&rarr;") {
			meanings = append(meanings, r)
		}
	}
	if len(meanings) > 10 {
		meanings = meanings[:10]
	}

	t := template.New("trans table")
	t.Parse(`
{{range .}}
# {{.Word}}
{{.Meaning}}
{{end}}
`)

	var tpl bytes.Buffer
	err = t.Execute(&tpl, meanings)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.String(), nil
}
