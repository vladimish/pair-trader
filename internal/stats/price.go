package stats

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type SellBuy struct {
	AvgPos float64 `json:"avgPos"`
	AvgNeg float64 `json:"avgNeg"`
}

func FindSellBuyPrices(data *sync.Map) *sync.Map {
	res := sync.Map{}

	data.Range(func(key any, value any) bool {
		avgNeg, avgPos := 0.0, 0.0
		negCnt, posCnt := 0, 0

		spread := value.([]float64)

		for j := 0; j < len(spread); j++ {
			if spread[j] < 0 {
				negCnt++
				avgNeg += spread[j]
			} else {
				posCnt++
				avgPos += spread[j]
			}
		}

		if posCnt != 0 {
			avgPos /= float64(posCnt)
		}
		if negCnt != 0 {
			avgNeg /= float64(negCnt)
		}

		logrus.Debug(key, avgPos, avgNeg)
		res.Store(key, SellBuy{AvgNeg: avgNeg, AvgPos: avgPos})

		return true
	})

	return &res
}
