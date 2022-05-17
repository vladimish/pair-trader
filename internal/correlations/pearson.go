package correlations

import (
	"github.com/montanaflynn/stats"
	"github.com/sirupsen/logrus"
	"github.com/vladimish/pair-trader/internal/data/models"
)

func BuildCorrelationMatrix(cd []models.CandlesData) [][]float64 {
	rs := make([][]float64, len(cd))

	for i := 0; i < len(cd)-1; i++ {
		rs[i] = make([]float64, len(cd))
		for j := i; j < len(cd); j++ {
			iPrices, jPrices := cd[i].ConvertPrices(), cd[j].ConvertPrices()
			r := FindPearsonCorrelation(iPrices, jPrices)
			//t += fmt.Sprintf("%f,", r)
			logrus.Debugf("correlation between %s and %s is %f", cd[i].Figi, cd[j].Figi, r)
			rs[i][j] = r
		}
	}

	return rs
}

func FindPearsonCorrelation(x, y []float64) float64 {
	res, err := stats.Correlation(x, y)
	if err != nil {
		panic(err)
	}

	return res
}
