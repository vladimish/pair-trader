package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/api"
	"github.com/vladimish/pair-trader/internal/debug"
	"github.com/vladimish/pair-trader/internal/env"
	"github.com/vladimish/pair-trader/internal/models"
	"github.com/vladimish/pair-trader/internal/stats"
	"io"
	"os"
	"sync"
	"time"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	//l, err := os.Create(time.Now().String() + ".log")
	//if err != nil {
	//	panic(err)
	//}
	//logrus.SetOutput(l)

	logrus.Info("starting bot...")
	err := env.InitEnv()
	if err != nil {
		logrus.Panic("can't connect to the api: ", err)
	}

	df, err := os.Open("data.json")
	var cd []models.CandlesData
	if err != nil {
		logrus.Info("fetching data...")
		figis, err := env.E.SDK.ShareTickersToFigis(env.E.CFG.Tickers)
		if err != nil {
			panic(err)
		}
		figis = append(figis, env.E.CFG.Figis...)

		cd, err = api.FetchDataAndComplete(figis, time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now(), investapi.CandleInterval_CANDLE_INTERVAL_HOUR)
		if err != nil {
			panic(err)
		}

		res, err := json.Marshal(cd)
		if err != nil {
			panic(err)
		}
		df, err := os.Create("data.json")
		if err != nil {
			panic(err)
		}
		_, err = df.Write(res)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		logrus.Info("loading data...")
		bytes, err := io.ReadAll(df)
		if err != nil {
			panic(err)
		}

		logrus.Info("marshaling data...")
		err = json.Unmarshal(bytes, &cd)
		if err != nil {
			panic(err)
		}
	}

	logrus.Info("building correlation matrix...")
	fc, err := os.Create("cor.csv")
	if err != nil {
		panic(err)
	}
	t := ","
	for i := range cd {
		t += fmt.Sprintf("%s,", cd[i].Figi)
	}
	_, err = fc.Write([]byte(t[:len(t)-1] + "\n"))
	if err != nil {
		logrus.Error(err)
	}

	rs := stats.BuildCorrelationMatrix(cd)

	for i := range rs {
		t := cd[i].Figi + ","
		for j := range rs[i] {
			t += fmt.Sprintf("%f,", rs[i][j])
		}
		_, err = fc.Write([]byte(t[:len(t)-1] + "\n"))
		if err != nil {
			logrus.Error(err)
		}
	}

	err = debug.BuildAndSavePricePlot(cd)
	if err != nil {
		logrus.Error(err)
	}
	err = debug.BuildAndSaveNormalizedPlot(cd)
	if err != nil {
		logrus.Error(err)
	}
	err = debug.BuildAndSavePercentagePlot(cd)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("building spread plots...")
	wg := sync.WaitGroup{}

	goodSpreads, goodTimes := sync.Map{}, sync.Map{}
	for i := 0; i < len(cd)-1; i++ {
		for j := i + 1; j < len(cd); j++ {
			wg.Add(1)
			go func(i, j int) {
				if rs[i][j] > 0.7 {
					name := cd[i].Figi + "-" + cd[j].Figi
					//fmt.Println(cd[i].Figi, cd[j].Figi, i, j)
					s, t := stats.BuildSpreadPlot(cd[i], cd[j])
					goodSpreads.Store(name, s)
					goodTimes.Store(name, t)

					err := debug.SaveSpread(s, t, name)
					if err != nil {
						logrus.Error(err)
					}
				}
				wg.Done()
			}(i, j)
		}
	}

	wg.Wait()

	goodSpreads.Range(func(key any, value any) bool {
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

		fmt.Println(key, avgPos, avgNeg)

		return true
	})

}
