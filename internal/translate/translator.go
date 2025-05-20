package translate

type TranslatorConfig struct {
	To       string
	LLMmodel string
	ApiKey   string
}

type Translator interface {
	Translate(text string) (string, error)
}

func New(t string, cfg TranslatorConfig) Translator {
	switch t {
	case "grok":
		return grok{To: cfg.To, LlmModel: cfg.LLMmodel, ApiKey: cfg.ApiKey}
	case "google":
		return google{To: cfg.To}
	}
	return google{To: cfg.To}
}
