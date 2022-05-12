package correlations

import (
	"github.com/sirupsen/logrus"
	"github.com/vladimish/pair-trader/internal/data/models"
	"math/big"
)

func BuildCorrelationMatrix(cd []models.CandlesData) [][]float64 {
	rs := make([][]float64, len(cd))

	for i := 0; i < len(cd)-1; i++ {
		rs[i] = make([]float64, len(cd))
		for j := i; j < len(cd); j++ {
			iPrices, jPrices := cd[i].ConvertPrices(), cd[j].ConvertPrices()
			r := FindPearsonCorrelation(iPrices, jPrices)
			//t += fmt.Sprintf("%f,", r)
			logrus.Infof("correlation between %s and %s is %f", cd[i].Figi, cd[j].Figi, r)
			rs[i][j] = r
		}
	}

	return rs
}

func FindPearsonCorrelation(x, y []float64) float64 {
	big.NewFloat(0)
	var xSum, ySum, xySum, xSqr, ySqr = make(chan *big.Float), make(chan *big.Float), make(chan *big.Float), make(chan *big.Float), make(chan *big.Float)

	go func(data []float64, res chan *big.Float) {
		t := sum(data)
		res <- t
	}(x, xSum)
	go func(data []float64, res chan *big.Float) {
		t := sum(data)
		res <- t
	}(y, ySum)
	go func(xarr, yarr []float64, res chan *big.Float) {
		t := mulSum(xarr, yarr)
		res <- t
	}(x, y, xySum)
	go func(data []float64, res chan *big.Float) {
		t := sqrSum(data)
		res <- t
	}(x, xSqr)
	go func(data []float64, res chan *big.Float) {
		t := sqrSum(data)
		res <- t
	}(y, ySqr)

	xr, yr, xyr, xsr, ysr := <-xSum, <-ySum, <-xySum, <-xSqr, <-ySqr
	//fmt.Println(xr, yr, xyr, xsr, ysr)

	divisor, dividenl, dividenr := big.NewFloat(0), big.NewFloat(0), big.NewFloat(0)
	divisor.Sub(big.NewFloat(0).Mul(big.NewFloat(float64(len(x))), xyr), big.NewFloat(0).Mul(xr, yr))
	dividenl = big.NewFloat(0).Sub(big.NewFloat(0).Mul(big.NewFloat(float64(len(x))), xsr), big.NewFloat(0).Mul(xr, xr))
	dividenr = big.NewFloat(0).Sub(big.NewFloat(0).Mul(big.NewFloat(float64(len(x))), ysr), big.NewFloat(0).Mul(yr, yr))

	r := big.NewFloat(0)
	r.Quo(divisor, big.NewFloat(0).Sqrt(big.NewFloat(0).Mul(dividenl, dividenr)))

	res, _ := r.Float64()
	//fmt.Println(res, acc)
	return res
}

func sum(arr []float64) (res *big.Float) {
	res = big.NewFloat(0)
	for _, el := range arr {
		res.Add(res, big.NewFloat(el))
	}

	return res
}

func mulSum(x, y []float64) (res *big.Float) {
	res = big.NewFloat(0)
	for i := range x {
		res.Add(res, big.NewFloat(x[i]*y[i]))
	}

	return res
}

func sqrSum(arr []float64) (res *big.Float) {
	res = big.NewFloat(0)
	for _, el := range arr {
		res.Add(res, big.NewFloat(el*el))
	}

	return res
}
