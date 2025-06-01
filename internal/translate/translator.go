package translate

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type TranslatorConfig struct {
	To       string
	LLMmodel string
	ApiKey   string
}

type TransType int

const (
	TypeGroq          TransType = iota
	TypeGroqAlayzer             = iota
	TypeGoogle                  = iota
	TypeDictionaryApi           = iota
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
		return groq{To: nativeTargetLang, LlmModel: cfg.LLMmodel, ApiKey: cfg.ApiKey}
	case TypeGoogle:
		return google{To: cfg.To}
	case TypeDictionaryApi:
		return dictionaryapi{}
	case TypeGroqAlayzer:
		return groqAnalyzer{LlmModel: cfg.LLMmodel, ApiKey: cfg.ApiKey}
	}
	return google{To: cfg.To}
}
