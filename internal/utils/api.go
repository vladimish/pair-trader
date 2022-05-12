package utils

import investapi "github.com/tinkoff/invest-api-go-sdk"

func IntervalToSeconds(i investapi.CandleInterval) int {
	switch i {
	case investapi.CandleInterval_CANDLE_INTERVAL_1_MIN:
		return 60
	case investapi.CandleInterval_CANDLE_INTERVAL_5_MIN:
		return 60 * 5
	case investapi.CandleInterval_CANDLE_INTERVAL_15_MIN:
		return 60 * 15
	case investapi.CandleInterval_CANDLE_INTERVAL_HOUR:
		return 60 * 60
	case investapi.CandleInterval_CANDLE_INTERVAL_DAY:
		return 60 * 60 * 24
	}

	return 0
}
