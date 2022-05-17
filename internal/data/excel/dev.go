package excel

import (
	"errors"
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
	f.SetCellValue(name, "A1", "time")
	f.SetCellValue(name, "B1", name)

	for i := 0; i < len(spread); i++ {
		f.SetCellValue(name, "A"+strconv.Itoa(i+2), time[i])
		f.SetCellValue(name, "B"+strconv.Itoa(i+2), spread[i])
	}

	err = f.SaveAs(SPREAD_PATH)
	m.Unlock()
	if err != nil {
		return err
	}

	return nil
}
