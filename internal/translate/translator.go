package translate

type TranslatorConfig struct {
	To       string
	LLMmodel string
	ApiKey   string
}

type TransType int

const (
	TypeGroq          TransType = iota
	TypeGoogle                  = iota
	TypeDictionaryApi           = iota
)

type Translator interface {
	Translate(text string) (string, error)
}

func New(t TransType, cfg TranslatorConfig) Translator {
	switch t {
	case TypeGroq:
		return groq{To: cfg.To, LlmModel: cfg.LLMmodel, ApiKey: cfg.ApiKey}
	case TypeGoogle:
		return google{To: cfg.To}
	case TypeDictionaryApi:
		return dictionaryapi{}
	}
	return google{To: cfg.To}
}
