package service

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogFile struct {
	Id             string
	Cmd            string
	StartTimestamp time.Time
	Output         string
}

func GetLogFiles(logdir string) ([]LogFile, error) {
	logs := []LogFile{}
	err := filepath.WalkDir(logdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == logdir {
				return nil
			}
			dirPath := filepath.Join(logdir, d.Name())
			dir := os.DirFS(dirPath)
			log.Default().Print(dir)
			f := LogFile{}
			if cmdFile, err := fs.Stat(dir, "cmd"); err == nil {
				f.Id = dirPath
				f.StartTimestamp = cmdFile.ModTime()
			}
			if cmd, err := fs.ReadFile(dir, "cmd"); err == nil {
				f.Cmd = string(cmd)
			} else {
				return err
			}
			if _, err := fs.Stat(dir, "output"); err == nil {
				f.Output = filepath.Join(dirPath, "output")
			} else {
				return nil
			}
			logs = append(logs, f)
		}
		return nil
	})
	return logs, err
}
