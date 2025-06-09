package translate

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type TranslatorConfig struct {
	To        string
	GroqModel string
	ApiKey    string
}

type TransType string

const (
	TypeGroq          = "groq"
	TypeGroqAlayzer   = "groq_analyze"
	TypeGoogle        = "google"
	TypeDictionaryApi = "dictapi"
)

type Translator interface {
	Translate(text string) (string, error)
}

func New(t TransType, cfg TranslatorConfig) Translator {
	nativeTargetLang := cfg.To
	l, err := language.Parse(cfg.To)
	if err == nil {
		nativeTargetLang = display.Self.Name(l)
	}

	switch t {
	case TypeGroq:
		return groq{To: nativeTargetLang, LlmModel: cfg.GroqModel, ApiKey: cfg.ApiKey}
	case TypeGoogle:
		return google{To: cfg.To}
	case TypeDictionaryApi:
		return dictionaryapi{}
	case TypeGroqAlayzer:
		return groqAnalyzer{LlmModel: cfg.GroqModel, ApiKey: cfg.ApiKey}
	}
	return google{To: cfg.To}
}
