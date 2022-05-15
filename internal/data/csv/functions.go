package csv

import (
	"fmt"
	"github.com/vladimish/pair-trader/internal/data/models"
	"os"
)

func BuildAndSavePricePlot(cd []models.CandlesData) {
	fn, err := os.Create("res.csv")
	if err != nil {
		panic(err)
	}
	prices := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		prices[i] = cd[i].ConvertPrices()
	}
	fn.Write([]byte(t[:len(t)-1] + "\n"))
	for i := 0; i < len(prices[0]); i++ {
		t := fmt.Sprintf("%d,", cd[0].Candles[i].Time.Seconds)
		for j := range prices {
			t += fmt.Sprintf("%f,", prices[j][i])
		}
		fn.Write([]byte(t[:len(t)-1] + "\n"))
	}
}

func BuildAndSaveNormalizedPlot(cd []models.CandlesData) {
	fn, err := os.Create("nres.csv")
	if err != nil {
		panic(err)
	}
	normalized := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		normalized[i] = cd[i].Normalize()
	}
	fn.Write([]byte(t[:len(t)-1] + "\n"))
	for i := 0; i < len(normalized[0]); i++ {
		t := fmt.Sprintf("%d,", cd[0].Candles[i].Time.Seconds)
		for j := range normalized {
			t += fmt.Sprintf("%f,", normalized[j][i])
		}
		fn.Write([]byte(t[:len(t)-1] + "\n"))
	}
}

func BuildAndSavePercentagePlot(cd []models.CandlesData) {
	fn, err := os.Create("pres.csv")
	if err != nil {
		panic(err)
	}
	percentage := make([][]float64, len(cd))
	t := ","
	for i := range cd {
		t += cd[i].Figi + ","
		percentage[i] = cd[i].Percent()
	}

	//for k := range percentage {
	//	percentage[k][0] = 0
	//	for i := 1; i < len(percentage[k]); i++ {
	//		percentage[k][i] = percentage[k][i-1] + percentage[k][i]
	//	}
	//}

	fn.Write([]byte(t[:len(t)-1] + "\n"))
	for i := 0; i < len(percentage[0]); i++ {
		t := fmt.Sprintf("%d,", cd[0].Candles[i].Time.Seconds)
		for j := range percentage {
			t += fmt.Sprintf("%f,", percentage[j][i])
		}
		fn.Write([]byte(t[:len(t)-1] + "\n"))
	}
}
