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

	logrus.Info("starting bot...")
	err := env.InitEnv()
	if err != nil {
		logrus.Panic("can't connect to the api: ", err)
	}

	df, err := os.Open("data.json")
	var cd []models.CandlesData
	if err != nil {
		logrus.Info("fetching data...")
		cd, err = correlations.FetchDataAndComplete(env.E.CFG.Figis, time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC), time.Now(), investapi.CandleInterval_CANDLE_INTERVAL_DAY)
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
		bytes, err := io.ReadAll(df)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &cd)
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
}
