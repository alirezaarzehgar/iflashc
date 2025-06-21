package config

import (
	"bytes"
	"text/template"

	_ "embed"

	"github.com/alirezaarzehgar/iflashc/internal/query"
)

type TransType string

const (
	TypeOpenAI        = "openai"
	TypeGoogle        = "google"
	TypeDictionaryApi = "dictapi"
	TypeFastdic       = "fastdic"
)

var ConfigurableTranslators = []string{
	TypeOpenAI,
	TypeGoogle,
	TypeDictionaryApi,
	TypeFastdic,
}

type Keys struct {
	DestLang      string
	Context       string
	Socks5        string
	Translator    string
	OpenAIApiKey  string
	OpenAIBaseURL string
	OpenAIModel   string
}

var (
	DefaultKeys = Keys{
		DestLang:      "dest_lang",
		Context:       "context",
		Socks5:        "proxy_socks5",
		Translator:    "translator",
		OpenAIApiKey:  "openai_api_key",
		OpenAIModel:   "openai_model",
		OpenAIBaseURL: "openai_base_url",
	}

	DefaultConfigs = Defaults{
		Translator: TypeDictionaryApi,
		DestLang:   "fa",
		Context:    "default",
	}

	ConfigurableKeys = []string{
		DefaultKeys.Translator,
		DefaultKeys.Context,
		DefaultKeys.OpenAIApiKey,
		DefaultKeys.OpenAIModel,
		DefaultKeys.OpenAIBaseURL,
		DefaultKeys.DestLang,
		DefaultKeys.Socks5,
	}
)

//go:embed schema.sql
var schemaTemplate string

type Defaults struct {
	Translator TransType
	DestLang   string
	Context    string
}

func GetSchema() (string, error) {
	schemaTemplate += `
INSERT INTO kvstore (key, value) VALUES 
('{{ .Keys.Translator }}', '{{ .Defaults.Translator }}'),
('{{ .Keys.DestLang }}', '{{ .Defaults.DestLang }}'),
('{{ .Keys.Context }}', '{{ .Defaults.Context }}');
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
		Defaults: DefaultConfigs,
		Keys:     DefaultKeys,
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
