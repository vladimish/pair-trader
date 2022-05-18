package debug

import (
	"errors"
	"fmt"
	"github.com/vladimish/pair-trader/internal/models"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
	"sync"
)

const (
	SPREAD_PATH = "spread.xlsx"
)

var m sync.Mutex

func SaveSpread(spread []float64, time []int64, name string) error {
	if len(spread) != len(time) {
		return errors.New("spread and time have different length")
	}

	m.Lock()
	defer m.Unlock()
	f, err := excelize.OpenFile(SPREAD_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			f = excelize.NewFile()
		} else {
			return err
		}
	}

	index := f.NewSheet(name)
	f.SetActiveSheet(index)
	err = f.SetCellValue(name, "A1", "time")
	if err != nil {
		return err
	}
	err = f.SetCellValue(name, "B1", name)
	if err != nil {
		return err
	}

	for i := 0; i < len(spread); i++ {
		err = f.SetCellValue(name, "A"+strconv.Itoa(i+2), time[i])
		if err != nil {
			return err
		}

		err = f.SetCellValue(name, "B"+strconv.Itoa(i+2), spread[i])
		if err != nil {
			return err
		}
	}

	err = f.SaveAs(SPREAD_PATH)
	if err != nil {
		return err
	}

	return nil
}

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
