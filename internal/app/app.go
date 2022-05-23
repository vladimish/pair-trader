package app

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/api"
	"github.com/vladimish/pair-trader/internal/env"
	"github.com/vladimish/pair-trader/internal/models"
	"github.com/vladimish/pair-trader/internal/stats"
	"io"
	"os"
	"sync"
	"time"
)

type App struct {
	Data []models.CandlesData
}

func NewApp() (*App, error) {
	cd, err := loadData(time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now())
	if err != nil {
		return nil, err
	}

	a := &App{
		Data: cd,
	}

	return a, nil
}

func (a *App) BuildCorrelationMatrix() ([][]float64, error) {
	fc, err := os.Create("cor.csv")
	if err != nil {
		return nil, err
	}

	t := ","
	for i := range a.Data {
		t += fmt.Sprintf("%s,", a.Data[i].Figi)
	}
	_, err = fc.Write([]byte(t[:len(t)-1] + "\n"))
	if err != nil {
		logrus.Error(err)
	}

	rs := stats.BuildCorrelationMatrix(a.Data)
	return rs, nil
}

func (a *App) BuildSpread(correlationMatrix [][]float64) (goodSpreads, goodTimes *sync.Map) {
	wg := sync.WaitGroup{}

	goodSpreads, goodTimes = &sync.Map{}, &sync.Map{}
	f, err := stats.NewSpreadFile("spread.xlsx")
	if err != nil {
		logrus.Error(err)
	}

	for i := 0; i < len(a.Data)-1; i++ {
		for j := i + 1; j < len(a.Data); j++ {
			wg.Add(1)
			go func(i, j int) {
				if correlationMatrix[i][j] > env.E.CFG.Params.MinCorrelation {
					name := a.Data[i].Figi + "-" + a.Data[j].Figi
					//fmt.Println(cd[i].Figi, cd[j].Figi, i, j)
					s, t := stats.BuildSpreadPlot(a.Data[i], a.Data[j])
					goodSpreads.Store(name, s)
					goodTimes.Store(name, t)

					err := f.AddSpread(s, t, name)
					if err != nil {
						logrus.Error(err)
					}
				}
				wg.Done()
			}(i, j)
		}
	}

	wg.Wait()
	err = f.SaveSpread()
	if err != nil {
		logrus.Error(err)
	}

	return goodSpreads, goodTimes
}

func loadData(from, to time.Time) ([]models.CandlesData, error) {
	df, err := os.Open("data.json")
	if err != nil {
		cd := make([]models.CandlesData, 0)

		logrus.Info("fetching data...")
		figis, err := env.E.SDK.ShareTickersToFigis(env.E.CFG.Tickers)
		if err != nil {
			return nil, err
		}

		figis = append(figis, env.E.CFG.Figis...)
		cd, err = api.FetchDataAndAddMissing(figis, from, to, investapi.CandleInterval_CANDLE_INTERVAL_DAY)
		if err != nil {
			return nil, err
		}

		res, err := json.Marshal(cd)
		if err != nil {
			return nil, err
		}

		df, err := os.Create("data.json")
		if err != nil {
			return nil, err
		}

		_, err = df.Write(res)
		if err != nil {
			return nil, err
		}

		return cd, nil

	} else if !os.IsNotExist(err) {
		cd, err := loadFromFile(df)
		if err != nil {
			return nil, err
		}

		return api.AddMissing(cd), nil
	} else {
		return nil, err
	}
}

func loadFromFile(df *os.File) ([]models.CandlesData, error) {
	var cd []models.CandlesData

	logrus.Info("loading data...")
	bytes, err := io.ReadAll(df)
	if err != nil {
		return nil, err
	}

	logrus.Info("marshaling data...")
	err = json.Unmarshal(bytes, &cd)
	if err != nil {
		return nil, err
	}

	return cd, nil
}
