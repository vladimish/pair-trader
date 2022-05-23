package realtime

import (
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/env"
	"strings"
	"sync"
)

func GetPairData(pairs *sync.Map) ([]*investapi.LastPrice, error) {
	figis := make([]string, 0)
	pairs.Range(func(key, value any) bool {
		split := strings.Split(key.(string), "-")
		figis = append(figis, split[0], split[1])
		return true
	})

	prices, err := env.E.SDK.GetLastPrices(figis)

	return prices, err
}
