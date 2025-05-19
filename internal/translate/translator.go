package translate

import "os"

type Translator interface {
	Translate(text string) (string, error)
}

type TranslatorType int

const (
	TypeGrok = iota
	TypeGoogleTranslate
)

func New(t TranslatorType, to string) Translator {
	switch t {
	case TypeGrok:
		return grok{To: to, LlmModel: os.Getenv("GROK_LLM_MODEL"), ApiKey: os.Getenv("GROK_API_KEY")}
	case TypeGoogleTranslate:
		return google{To: to}
	}
	return grok{To: to}
}
