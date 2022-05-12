package env

import (
	"github.com/vladimish/pair-trader/pkg/config"
	"github.com/vladimish/pair-trader/pkg/invest"
)

type Env struct {
	SDK *invest.SDK
	CFG *config.Config
}

var E *Env

func InitEnv() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	sdk, err := invest.NewSDK(cfg.Token)
	if err != nil {
		return err
	}
	E = &Env{
		SDK: sdk,
		CFG: cfg,
	}

	return nil
}
