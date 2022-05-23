package stats

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"os"
	"strconv"
	"sync"
)

type SpreadFile struct {
	filePath string
	f        *excelize.File
	m        sync.Mutex
}

func NewSpreadFile(path string) (*SpreadFile, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			f = excelize.NewFile()
		} else {
			return nil, err
		}
	}

	s := &SpreadFile{
		filePath: path,
		f:        f,
		m:        sync.Mutex{},
	}

	return s, nil
}

func (s *SpreadFile) AddSpread(spread []float64, time []int64, name string) error {
	if len(spread) != len(time) {
		return errors.New("spread and time have different length")
	}

	s.m.Lock()
	defer s.m.Unlock()

	index := s.f.NewSheet(name)
	s.f.SetActiveSheet(index)
	err := s.f.SetCellValue(name, "A1", "time")
	if err != nil {
		return err
	}
	err = s.f.SetCellValue(name, "B1", name)
	if err != nil {
		return err
	}

	for i := 0; i < len(spread); i++ {
		err = s.f.SetCellValue(name, "A"+strconv.Itoa(i+2), time[i])
		if err != nil {
			return err
		}

		err = s.f.SetCellValue(name, "B"+strconv.Itoa(i+2), spread[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SpreadFile) SaveSpread() error {
	logrus.Debug("saving spread file...")
	return s.f.SaveAs(s.filePath)
}
