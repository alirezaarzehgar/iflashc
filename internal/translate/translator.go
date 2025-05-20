package translate

import "os"

type Translator interface {
	Translate(text string) (string, error)
}

func New(t string, to string) Translator {
	switch t {
	case "grok":
		return grok{To: to, LlmModel: os.Getenv("GROK_LLM_MODEL"), ApiKey: os.Getenv("GROK_API_KEY")}
	case "google":
		return google{To: to}
	}
	return google{To: to}
}
