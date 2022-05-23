package main

import (
	"github.com/sirupsen/logrus"
	"github.com/vladimish/pair-trader/internal/app"
	"github.com/vladimish/pair-trader/internal/env"
	"github.com/vladimish/pair-trader/internal/realtime"
	"github.com/vladimish/pair-trader/internal/stats"
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

	a, err := app.NewApp()
	if err != nil {
		logrus.Panic(err)
	}

	logrus.Info("building correlation matrix...")
	correlationMatrix, err := a.BuildCorrelationMatrix()
	if err != nil {
		logrus.Error(err)
	}

	err = stats.SaveCorrelationPlot(correlationMatrix, a.Data)
	if err != nil {
		logrus.Error(err)
	}

	err = stats.BuildAndSavePricePlot(a.Data)
	if err != nil {
		logrus.Error(err)
	}
	err = stats.BuildAndSaveNormalizedPlot(a.Data)
	if err != nil {
		logrus.Error(err)
	}
	err = stats.BuildAndSavePercentagePlot(a.Data)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("building spread plots...")
	gs, _ := a.BuildSpread(correlationMatrix)

	sellBuy := stats.FindSellBuyPrices(gs)

	for i := 0; i < 5; i++ {
		data, err := realtime.GetPairData(sellBuy)
		if err != nil {
			panic(err)
		}

		logrus.Debug(data)

		time.Sleep(5 * time.Second)
	}
}
