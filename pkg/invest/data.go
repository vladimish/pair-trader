package invest

import (
	"github.com/sirupsen/logrus"
	investapi "github.com/tinkoff/invest-api-go-sdk"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func getStep(interval investapi.CandleInterval) (step time.Duration) {
	switch interval {
	case investapi.CandleInterval_CANDLE_INTERVAL_1_MIN, investapi.CandleInterval_CANDLE_INTERVAL_5_MIN, investapi.CandleInterval_CANDLE_INTERVAL_15_MIN:
		step = time.Hour * 24
		break
	case investapi.CandleInterval_CANDLE_INTERVAL_HOUR:
		step = time.Hour * 24 * 7
		break
	case investapi.CandleInterval_CANDLE_INTERVAL_DAY:
		step = time.Hour * 24 * 365
	}

	return step
}

func (s SDK) GetCandles(figi string, from, to time.Time, interval investapi.CandleInterval) ([]*investapi.HistoricCandle, error) {
	step := getStep(interval)
	if !from.Add(step).Before(to) {
		step = to.Sub(from) - time.Second
	}

	candles := make([]*investapi.HistoricCandle, 0)
	for i := from.Add(step); i.Before(to); i = i.Add(step) {
		j := i.Add(-step)

		logrus.Debugf("sending GetCandles request from %s to %s...", j.String(), i.String())

		req := &investapi.GetCandlesRequest{
			Figi:     figi,
			From:     timestamppb.New(j),
			To:       timestamppb.New(i),
			Interval: interval,
		}

		var resp *investapi.GetCandlesResponse
		// Cycle is used to repeat request if ResourceExhausted got.
		for i := 0; i < 1; i++ {
			var err error
			resp, err = s.MarketData.GetCandles(s.ctx, req)
			if err != nil {
				if err.Error() == "rpc error: code = ResourceExhausted desc = " {
					logrus.Debug("exhausted, waiting 5 seconds...")
					time.Sleep(time.Second * 5)
					i--
					continue
				}
				return nil, err
			}
		}

		logrus.Debugf("got %d candles of %s", len(resp.GetCandles()), figi)

		candles = append(candles, resp.GetCandles()...)

		if !i.Add(step).Before(to) {
			step = to.Sub(i) - time.Second
			if step == 0 {
				break
			}
			if step < time.Hour*24 {
				break
			}
		}
	}

	return candles, nil
}

func (s *SDK) ShareTickersToFigis(tickers []string) ([]string, error) {
	res := make([]string, 0)

	resp, err := s.Instruments.Shares(s.ctx, &investapi.InstrumentsRequest{
		InstrumentStatus: investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE,
	})
	if err != nil {
		return nil, err
	}

T:
	for i := range tickers {
		for j := range resp.Instruments {
			if resp.Instruments[j].Ticker == tickers[i] {
				res = append(res, resp.Instruments[j].Figi)
				continue T
			}
		}
		logrus.Warn("ticker ", tickers[i], " not found")
	}

	return res, err
}

func (s *SDK) GetLastPrices(figis []string) ([]*investapi.LastPrice, error) {
	res, err := s.MarketData.GetLastPrices(s.ctx, &investapi.GetLastPricesRequest{
		Figi: figis,
	})
	if err != nil {
		return nil, err
	}

	return res.GetLastPrices(), nil
}
