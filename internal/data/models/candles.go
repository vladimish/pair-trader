package models

import (
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"github.com/vladimish/pair-trader/internal/utils"
	"math"
)

type CandlesData struct {
	Figi     string                      `json:"figi"`
	Interval int                         `json:"interval"`
	Candles  []*investapi.HistoricCandle `json:"candles"`
}

func (d *CandlesData) ConvertPrices() []float64 {
	res := make([]float64, len(d.Candles))
	for i := range res {
		l := utils.CountDigits(d.Candles[i].Close.Nano)
		dec := float64(d.Candles[i].Close.Nano) / math.Pow(10, float64(l))
		res[i] = float64(d.Candles[i].Close.Units) + dec
	}

	return res
}

func (d *CandlesData) Normalize() []float64 {
	data := make([]float64, len(d.Candles))
	data = d.ConvertPrices()

	min, max := math.MaxFloat64, 0.0
	for i := range data {
		if data[i] < min {
			min = data[i]
		}
		if data[i] > max {
			max = data[i]
		}
	}
	res := make([]float64, len(data))
	for i := range res {
		res[i] = (data[i] - min) / (max - min)
	}

	return res
}

func (d *CandlesData) Percent() []float64 {
	data := d.ConvertPrices()
	res := make([]float64, len(data)-1)

	for i := 1; i < len(data); i++ {
		if data[i]-data[i-1] > 0 {
			inc := data[i] - data[i-1]
			res[i-1] = inc / data[i-1]
		} else {
			dec := data[i-1] - data[i]
			res[i-1] = -dec / data[i-1]
		}
	}

	return res
}
