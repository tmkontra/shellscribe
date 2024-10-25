package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"time"

	"github.com/nxadm/tail"
)

type LogFile struct {
	Cmd            string
	StartTimestamp time.Time
	Output         string
}

func main() {
	usr, _ := user.Current()
	logdir := filepath.Join(usr.HomeDir, "mutter")
	logs, err := getLogFiles(logdir)
	if err != nil {
		log.Fatal(err)
	}
	slices.SortFunc(logs, func(a LogFile, b LogFile) int {
		return a.StartTimestamp.Compare(b.StartTimestamp)
	})
	for i, l := range logs {
		fmt.Printf("%d: %s\n", i+1, l.Cmd)
	}
	var i int
	fmt.Scanln(&i)
	l := logs[i-1]
	t, _ := tail.TailFile(l.Output, tail.Config{Follow: true})
	for line := range t.Lines {
		fmt.Println(line.Text)
	}
}

func getLogFiles(logdir string) ([]LogFile, error) {
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
