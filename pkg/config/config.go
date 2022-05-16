package config

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	Token   string   `json:"-"`
	Figis   []string `json:"figis"`
	Tickers []string `json:"tickers"`
}

func NewConfig() (*Config, error) {
	CFG := &Config{}

	path := flag.String("c", "config.json", "path to config")
	token := flag.String("t", "-", "invest api token")
	flag.Parse()

	if *token == "-" {
		panic("token is not set")
	}
	CFG.Token = *token

	bytes, err := os.ReadFile(*path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &CFG)
	if err != nil {
		return nil, err
	}

	return CFG, err
}
