package translate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"
)

const (
	MESSAGE_ROLE_USER      = "user"
	MESSAGE_ROLE_ASSISTANT = "assistant"

	GORQ_REQ_URL = "https://api.groq.com/openai/v1/chat/completions"
)

type ResponseFormat struct {
	Type string `json:"type"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GorqRequest struct {
	Messages       []Message       `json:"messages"`
	Model          string          `json:"model"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Temperature    float64         `json:"temperature,omitempty"`
	TopP           int             `json:"top_p,omitempty"`
	Stream         bool            `json:"stream,omitempty"`
	Stop           string          `json:"stop,omitempty"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

type GorqResponseError struct {
	Message         string `json:"message"`
	Type            string `json:"type"`
	Code            string `json:"code"`
	FailedGenerated string `json:"failed_generation"`
}

type GorqResponse struct {
	Error             GorqResponseError `json:"error"`
	ID                string            `json:"id"`
	Object            string            `json:"object"`
	Created           int64             `json:"created"`
	Model             string            `json:"model"`
	SystemFingerprint string            `json:"system_fingerprint"`
	Choices           []Choice          `json:"choices"`
}

type Groq struct {
}

func talkToGroq(prompt, llmModel, apiKey string) (*GorqResponse, error) {
	body := GorqRequest{
		Messages: []Message{{
			Role:    MESSAGE_ROLE_USER,
			Content: prompt,
		}},
		Model:          llmModel,
		Temperature:    1,
		MaxTokens:      1024,
		TopP:           1,
		ResponseFormat: &ResponseFormat{Type: "json_object"},
	}
	commitMessage, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, GORQ_REQ_URL, bytes.NewBuffer(commitMessage))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", req.Host, err)
	}
	resBytes, _ := io.ReadAll(res.Body)

	gorqRes := GorqResponse{}
	err = json.Unmarshal(resBytes, &gorqRes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d: %s", res.StatusCode, gorqRes.Error.Message)
	}

	return &gorqRes, nil
}

type TranslationResponse struct {
	Definitions []string `json:"definitions"`
	Synonyms    []string `json:"synonyms"`
	Examples    []string `json:"examples"`
	Meanings    []string `json:"meanings"`
}

type groq struct {
	To       string
	LlmModel string
	ApiKey   string
}

func (g groq) Translate(text string) (string, error) {
	prompt := `Translate and provide detailed information for the term/phrase <<%s>> into %s.  

The output must be a well-structured JSON object with the following fields:  
1. **Definitions**: Clear and concise English definitions of the term/phrase.  
2. **Synonyms**: A list of relevant English synonyms (if applicable).  
3. **Examples**: Example sentences in English demonstrating proper usage.  
4. **Meanings**: Accurate translations and meanings in the target language (%s).

### Output Format (Strict JSON):    
{  
  "definitions": ["definition_1", "definition_2", ...],  
  "synonyms": ["synonym_1", "synonym_2", ...],  
  "examples": ["example_sentence_1", "example_sentence_2", ...],  
  "meanings": ["translation_1_in_target_language", "translation_2_in_target_language", ...]  
}
`
	resp, err := talkToGroq(fmt.Sprintf(prompt, text, g.To, g.To), g.LlmModel, g.ApiKey)
	if err != nil {
		return "", err
	}

	tr := TranslationResponse{}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &tr)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal content to valid json: %w", err)
	}

	t := template.New("trans table")
	t.Parse(`
# Meanings
{{range .Meanings}}
- {{.}}
{{end}}

# Definitions
{{range .Definitions}}
- {{.}}
{{end}}

# Synonyms
{{range .Synonyms}}
- {{.}}
{{end}}

# Examples
{{range .Examples}}
- {{.}}
{{end}}
`)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, tr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.String(), nil
}
