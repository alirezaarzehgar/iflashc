package config

import (
	"bytes"
	"text/template"

	_ "embed"

	"github.com/alirezaarzehgar/iflashc/internal/query"
)

type Keys struct {
	DestLang      string
	Context       string
	Socks5        string
	Translator    string
	GroqApiKey    string
	GroqModel     string
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
		GroqApiKey:    "groq_api_key",
		GroqModel:     "groq_model",
		OpenAIApiKey:  "openai_api_key",
		OpenAIModel:   "openai_model",
		OpenAIBaseURL: "openai_base_url",
	}

	ConfigEntries = []string{
		DefaultKeys.Translator,
		DefaultKeys.Context,
		DefaultKeys.OpenAIApiKey,
		DefaultKeys.OpenAIModel,
		DefaultKeys.OpenAIBaseURL,
		DefaultKeys.DestLang,

		DefaultKeys.GroqApiKey,
		DefaultKeys.GroqModel,
		DefaultKeys.Socks5,
	}

	DefaultConfigs = Defaults{
		Translator: TypeDictionaryApi,
		DestLang:   "fa",
	}
)

//go:embed schema.sql
var schemaTemplate string

type TransType string

const (
	TypeOpenAI        = "openai"
	TypeGroq          = "groq"
	TypeGroqAlayzer   = "groq_analyze"
	TypeGoogle        = "google"
	TypeDictionaryApi = "dictapi"
)

type Defaults struct {
	Translator TransType
	DestLang   string
}

func GetSchema() (string, error) {
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
