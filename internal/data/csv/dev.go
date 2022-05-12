package csv

import (
	"fmt"
	"github.com/vladimish/pair-trader/internal/data/models"
	"os"
)

func SaveTime(cd []models.CandlesData, path string) error {
	f, _ := os.Create(path)

	for i := range cd {
		_, err := f.Write([]byte(fmt.Sprintf("%s,", cd[i].Figi)))
		if err != nil {
			return err
		}
	}
	_, err := f.Write([]byte("\n"))
	if err != nil {
		return err
	}

	for i := range cd[0].Candles {
		t := ""
		for j := range cd {
			t += fmt.Sprintf("%d,", cd[j].Candles[i].Time.Seconds)
		}
		_, err := f.Write([]byte(t[:len(t)-1] + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
