package translate

import (
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type TranslatorConfig struct {
	To            string
	GroqModel     string
	OpenAIModel   string
	OpenAIBaseURL string
	ApiKey        string
}

type Translator interface {
	Translate(text string) (string, error)
}

func New(t config.TransType, cfg config.Config) Translator {
	switch t {
	case config.TypeOpenAI:
		return openAI{
			To:      cfg[config.DefaultKeys.DestLang],
			Model:   cfg[config.DefaultKeys.OpenAIModel],
			ApiKey:  cfg[config.DefaultKeys.OpenAIApiKey],
			BaseURL: cfg[config.DefaultKeys.OpenAIBaseURL],
		}
	case config.TypeGroq:
		return groq{
			To:       cfg[config.DefaultKeys.DestLang],
			LlmModel: cfg[config.DefaultKeys.GroqModel],
			ApiKey:   cfg[config.DefaultKeys.GroqApiKey],
		}
	case config.TypeGoogle:
		return google{To: cfg[config.DefaultKeys.DestLang]}
	case config.TypeDictionaryApi:
		return dictionaryapi{}
	case config.TypeGroqAlayzer:
		return groq{
			To:       cfg[config.DefaultKeys.DestLang],
			LlmModel: cfg[config.DefaultKeys.GroqModel],
			ApiKey:   cfg[config.DefaultKeys.GroqApiKey],
		}
	}

	nativeTargetLang := cfg[config.DefaultKeys.DestLang]
	l, err := language.Parse(nativeTargetLang)
	if err == nil {
		nativeTargetLang = display.Self.Name(l)
	}
	return google{To: nativeTargetLang}
}
