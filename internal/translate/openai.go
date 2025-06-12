package translate

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type openAI struct {
	To      string
	Model   string
	ApiKey  string
	BaseURL string
}

func (oAi openAI) Translate(text string) (string, error) {
	prompt := `Translate and provide detailed information for the term/phrase <<%[1]v>> into %[2]v.

The output must be a well-structured Markdown file with the following fields:  
1. **Definitions**: Clear and concise English definitions of the term/phrase.  
2. **Synonyms**: A list of relevant English synonyms (if applicable).  
3. **Examples**: Example sentences in English demonstrating proper usage.  
4. **Meanings**: Accurate translations and meanings in the target language (%[2]v).

Output Format (Strict Markdown):
# Meanings
- meaning1
- meaning2
...

# Definitions
- definition1
- definition2
...

# Synonyms
- synonym1
- synonym2
...

# Examples
- example1
- example2
...
`

	cfg := openai.DefaultConfig(oAi.ApiKey)
	cfg.BaseURL = oAi.BaseURL
	client := openai.NewClientWithConfig(cfg)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: oAi.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf(prompt, text, oAi.To),
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("ChatCompletion error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
