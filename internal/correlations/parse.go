package correlations

import (
	"github.com/sirupsen/logrus"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/data/models"
	"github.com/vladimish/pair-trader/internal/env"
	"github.com/vladimish/pair-trader/internal/utils"
	"golang.org/x/exp/slices"
	"time"
)

func FetchDataAndComplete(figis []string, from, to time.Time, interval investapi.CandleInterval) ([]models.CandlesData, error) {
	logrus.Info("getting historic candles...")
	d, err := FetchData(figis, from, to, interval)
	if err != nil {
		return nil, err
	}
	logrus.Info("completing data...")
	AddMissing(d)
	return d, nil
}

func FetchData(figis []string, from, to time.Time, interval investapi.CandleInterval) ([]models.CandlesData, error) {
	cd := make([]models.CandlesData, len(env.E.CFG.Figis))
	for i, s := range figis {
		logrus.Infof("getting %s candles...", s)
		candles, err := env.E.SDK.GetCandles(s, from, to, interval)
		if err != nil {
			logrus.Error("error while getting candles history: ", err)
		}

		cd[i] = models.CandlesData{
			Figi:     s,
			Interval: utils.IntervalToSeconds(interval),
			Candles:  candles,
		}
		//candlesMap[s] = candles
	}

	return cd, nil
}

func AddMissing(cd []models.CandlesData) {
	for i := 0; i < len(cd)-1; i++ {
		for j := i; j < len(cd); j++ {
			if j == i {
				continue
			}
		K:
			for k := 1; k < len(cd[i].Candles); k++ {
				if k == len(cd[j].Candles) {
					for ; k < len(cd[i].Candles); k++ {
						cd[j].Candles = insert(cd[j].Candles, cd[i].Candles, k)
					}
					break K
				}
				if cd[i].Candles[k].Time.Seconds < cd[j].Candles[k].Time.Seconds {
					cd[j].Candles = insert(cd[j].Candles, cd[i].Candles, k)
					k--
				} else if cd[i].Candles[k].Time.Seconds > cd[j].Candles[k].Time.Seconds {
					cd[i].Candles = insert(cd[i].Candles, cd[j].Candles, k)
					k--
				}
			}
		}
	}

	// Get rid of the first candle of each stock
	// because it was only used to add missing data.
	for i := range cd {
		//fmt.Println(len(cd[i].Candles))
		cd[i].Candles = cd[i].Candles[1:]
	}

}

func insert(candles []*investapi.HistoricCandle, timeCandles []*investapi.HistoricCandle, k int) []*investapi.HistoricCandle {
	return slices.Insert(candles, k, &investapi.HistoricCandle{
		Open:       candles[k-1].GetClose(),
		High:       candles[k-1].GetClose(),
		Low:        candles[k-1].GetClose(),
		Close:      candles[k-1].GetClose(),
		Volume:     0,
		Time:       timeCandles[k].GetTime(),
		IsComplete: candles[k-1].GetIsComplete(),
	})
}
