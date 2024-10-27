package service

import (
	"context"
	"log"
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
		return -a.StartTimestamp.Compare(b.StartTimestamp)
	})
	return files, nil
}

func (s *Service) TailFile(path string) (func(context.Context, chan string), error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return func(ctx context.Context, c chan string) {
		t, err := tail.TailFile(path, tail.Config{Follow: true, Poll: true})
		if err != nil {
			log.Printf("tail error: %s\n", err)
			close(c)
		}
		for {
			select {
			case <-ctx.Done():
				close(c)
				return
			case line := <-t.Lines:
				c <- line.Text
			}
		}
	}, nil
}
