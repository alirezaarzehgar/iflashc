package translate

import (
	"github.com/alirezaarzehgar/iflashc/internal/config"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

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
	case config.TypeGoogle:
		return google{To: cfg[config.DefaultKeys.DestLang]}
	case config.TypeDictionaryApi:
		return dictionaryapi{}
	case config.TypeFastdic:
		return fastdic{}
	}

	nativeTargetLang := cfg[config.DefaultKeys.DestLang]
	l, err := language.Parse(nativeTargetLang)
	if err == nil {
		nativeTargetLang = display.Self.Name(l)
	}
	return google{To: nativeTargetLang}
}
