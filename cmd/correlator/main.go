package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/correlations"
	"github.com/vladimish/pair-trader/internal/data/csv"
	"github.com/vladimish/pair-trader/internal/data/models"
	"github.com/vladimish/pair-trader/internal/env"
	"io"
	"os"
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

		cd, err = correlations.FetchDataAndComplete(figis, time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now(), investapi.CandleInterval_CANDLE_INTERVAL_DAY)
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
		df.Write(res)
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
	fc.Write([]byte(t[:len(t)-1] + "\n"))

	rs := correlations.BuildCorrelationMatrix(cd)

	for i := range rs {
		t := cd[i].Figi + ","
		for j := range rs[i] {
			t += fmt.Sprintf("%f,", rs[i][j])
		}
		fc.Write([]byte(t[:len(t)-1] + "\n"))
	}

	csv.BuildAndSavePricePlot(cd)
	csv.BuildAndSaveNormalizedPlot(cd)
	csv.BuildAndSavePercentagePlot(cd)

	for i := 0; i < len(cd)-1; i++ {
		for j := i + 1; j < len(cd); j++ {
			csv.BuildAndSaveSpreadPlot(cd, i, j, cd[i].Figi+"-"+cd[j].Figi+".csv")
		}
	}
}
