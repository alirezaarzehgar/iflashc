package translate

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
		return grok{To: to}
	case TypeGoogleTranslate:
		return google{To: to}
	}
	return grok{To: to}
}
