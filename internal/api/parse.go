package api

import (
	"github.com/sirupsen/logrus"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/env"
	"github.com/vladimish/pair-trader/internal/models"
	"github.com/vladimish/pair-trader/internal/utils"
	"golang.org/x/exp/slices"
	"time"
)

func FetchDataAndAddMissing(figis []string, from, to time.Time, interval investapi.CandleInterval) ([]models.CandlesData, error) {
	logrus.Info("getting historic candles...")
	d, err := FetchData(figis, from, to, interval)
	if err != nil {
		return nil, err
	}
	logrus.Info("completing data...")
	d = AddMissing(d)
	return d, nil
}

func FetchData(figis []string, from, to time.Time, interval investapi.CandleInterval) ([]models.CandlesData, error) {
	cd := make([]models.CandlesData, len(figis))
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

func AddMissing(cd []models.CandlesData) []models.CandlesData {
	for i := 0; i < len(cd); i++ {
		if len(cd[i].Candles) == 0 {
			cd = append(cd[:i], cd[i+1:]...)
			i--
		}
	}

	for i := 0; i < len(cd)-1; i++ {
		for j := i; j < len(cd); j++ {
			if j == i {
				continue
			}
		K:
			for k := 0; k < len(cd[i].Candles); k++ {
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
	for i := 0; i < len(cd); i++ {
		//fmt.Println(len(cd[i].Candles))
		cd[i].Candles = cd[i].Candles[1:]
	}

	return cd
}

func insert(candles []*investapi.HistoricCandle, timeCandles []*investapi.HistoricCandle, k int) []*investapi.HistoricCandle {
	knownIndex := k - 1
	if k == 0 {
		knownIndex = 1
	}

	return slices.Insert(candles, k, &investapi.HistoricCandle{
		Open:       candles[knownIndex].GetClose(),
		High:       candles[knownIndex].GetClose(),
		Low:        candles[knownIndex].GetClose(),
		Close:      candles[knownIndex].GetClose(),
		Volume:     0,
		Time:       timeCandles[k].GetTime(),
		IsComplete: candles[knownIndex].GetIsComplete(),
	})
}
