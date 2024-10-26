package service

import (
	"os"
	"slices"

	"github.com/nxadm/tail"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListFiles(dir string) ([]LogFile, error) {
	files, err := GetLogFiles(dir)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(files, func(a LogFile, b LogFile) int {
		return a.StartTimestamp.Compare(b.StartTimestamp)
	})
	return files, nil
}

func (s *Service) TailFile(path string) (func(chan string), error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return func(c chan string) {
		t, _ := tail.TailFile(path, tail.Config{Follow: true, Poll: true})
		for line := range t.Lines {
			c <- line.Text
		}
	}, nil
}
