package translate

import (
	"encoding/json"
	"fmt"
)

type groqAnalyzeResponse struct {
	Explanation string `json:"exp"`
}

type groqAnalyzer struct {
	LlmModel string
	ApiKey   string
}

// NOTE: Temporarily removed RTL (Right-to-Left) language support due to Fyne's current lack of
// native RTL text handling. Implementing a reliable bidirectional text analyzer within Fyne's
// limitations proved unfeasible. Re-evaluate when Fyne adds proper RTL support in future releases.
func (g groqAnalyzer) Translate(text string) (string, error) {
	prompt := `Analyze the given text and provide a concise, easy-to-understand explanation in English.

### Instructions:  
- Keep explanations simple, clear, and jargon-free. Aim for a 5th-grade reading level.  
- Output must be strict JSON format with no deviations.  

### Output Format (Strict JSON):  
{
  "exp": "Extremely simplified explanation in English."
}

### Given Text:  
%[1]v
`

	resp, err := talkToGroq(fmt.Sprintf(prompt, text), g.LlmModel, g.ApiKey)
	if err != nil {
		return "", err
	}

	tr := groqAnalyzeResponse{}
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &tr)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal content to valid json: %w", err)
	}

	return tr.Explanation, nil
}
