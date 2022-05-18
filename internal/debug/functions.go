package debug

import (
	"fmt"
	"github.com/vladimish/pair-trader/internal/models"
	"os"
)

func BuildAndSavePricePlot(cd []models.CandlesData) error {
	prices := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		prices[i] = cd[i].ConvertPrices()
	}

	return WriteCSV(cd, prices, t, "res.csv")
}

func BuildAndSaveNormalizedPlot(cd []models.CandlesData) error {
	normalized := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		normalized[i] = cd[i].Normalize()
	}

	return WriteCSV(cd, normalized, t, "nres.csv")
}

func BuildAndSavePercentagePlot(cd []models.CandlesData) error {
	percentage := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		percentage[i] = cd[i].Percent()
	}

	return WriteCSV(cd, percentage, t, "pres.csv")
}

func WriteCSV(cd []models.CandlesData, data [][]float64, t, filename string) error {
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