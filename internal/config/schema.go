package config

import (
	"bytes"
	"text/template"

	_ "embed"

	"github.com/alirezaarzehgar/iflashc/internal/query"
	"github.com/alirezaarzehgar/iflashc/internal/translate"
)

type Keys struct {
	DestLang   string
	Socks5     string
	Translator string
	GroqApiKey string
	GroqModel  string
}

var (
	defaultKeys = Keys{
		DestLang:   "dest_lang",
		Socks5:     "proxy_socks5",
		Translator: "translator",
		GroqApiKey: "groq_api_key",
		GroqModel:  "groq_model",
	}
)

//go:embed schema.sql
var schemaTemplate string

type Defaults struct {
	Translator translate.TransType
	DestLang   string
}

func GetSchema(vals Defaults) (string, error) {
	schemaTemplate += `
INSERT INTO kvstore (key, value) VALUES ('{{ .Keys.Translator }}', '{{ .Defaults.Translator }}');
INSERT INTO kvstore (key, value) VALUES ('{{ .Keys.DestLang }}', '{{ .Defaults.DestLang }}');
`

	t := template.New("Schema")
	tmpl, err := t.Parse(schemaTemplate)
	if err != nil {
		return "", err
	}

	schemaConsts := struct {
		Defaults Defaults
		Keys     Keys
	}{
		Defaults: vals,
		Keys:     defaultKeys,
	}

	res := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(res, schemaConsts); err != nil {
		return "", err
	}

	return res.String(), nil
}

type Config map[string]string

func ConfigToMap(kvs []query.Kvstore) Config {
	parsedConfig := Config{}
	for _, conf := range kvs {
		parsedConfig[conf.Key] = conf.Value
	}
	return parsedConfig
}
