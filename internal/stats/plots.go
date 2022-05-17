package stats

import "github.com/vladimish/pair-trader/internal/data/models"

func BuildSpreadPlot(x, y models.CandlesData) (spread []float64, time []int64) {
	if len(x.Candles) != len(y.Candles) {
		panic("x and y have different length")
	}

	spread = make([]float64, len(x.Candles))
	time = make([]int64, len(x.Candles))

	d1 := x.Percent()
	d2 := y.Percent()
	for k := range d1 {
		spread[k] = d1[k] - d2[k]
		time[k] = x.Candles[k].Time.Seconds
	}

	return spread, time
}
