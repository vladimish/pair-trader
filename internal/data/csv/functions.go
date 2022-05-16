package csv

import (
	"fmt"
	"github.com/vladimish/pair-trader/internal/data/models"
	"os"
)

func BuildAndSavePricePlot(cd []models.CandlesData) {
	prices := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		prices[i] = cd[i].ConvertPrices()
	}

	writeCsv(cd, prices, t, "res.csv")
}

func BuildAndSaveNormalizedPlot(cd []models.CandlesData) {
	normalized := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		normalized[i] = cd[i].Normalize()
	}

	writeCsv(cd, normalized, t, "nres.csv")
}

func BuildAndSavePercentagePlot(cd []models.CandlesData) {
	percentage := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		percentage[i] = cd[i].Percent()
	}

	writeCsv(cd, percentage, t, "pres.csv")
}

func BuildAndSaveSpreadPlot(cd []models.CandlesData, i, j int, filename string) {
	spread := make([]float64, len(cd[i].Candles))
	t := "," + cd[i].Figi + "-" + cd[j].Figi + ","

	d1 := cd[i].Percent()
	d2 := cd[j].Percent()
	for k := range cd[i].Candles {
		spread[k] = d1[k] - d2[k]
	}

	writeCsv(cd, [][]float64{spread}, t, filename)
}

func writeCsv(cd []models.CandlesData, data [][]float64, t, filename string) error {
	fn, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = fn.Write([]byte(t[:len(t)-1] + "\n"))
	if err != nil {
		return err
	}

	for i := 0; i < len(data[0]); i++ {
		t = fmt.Sprintf("%d,", cd[0].Candles[i].Time.Seconds)
		for j := range data {
			t += fmt.Sprintf("%f,", data[j][i])
		}
		_, err = fn.Write([]byte(t[:len(t)-1] + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
