package translate

import (
	gtranslate "github.com/gilang-as/google-translate"
)

type google struct {
	To string
}

func (g google) Translate(text string) (string, error) {
	v := gtranslate.Translate{
		Text: text,
		To:   g.To,
	}
	t, err := gtranslate.Translator(v)
	if err != nil {
		return "", err
	}

	response := "## " + t.Text + "\n"
	if t.Pronunciation != nil {
		response += *t.Pronunciation + "\n"
	}

	return response, nil
}
