package main

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/nxadm/tail"
	"github.com/tmkontra/shellscribe/internal/server"
	"github.com/tmkontra/shellscribe/internal/service"
)

func main() {
	usr, _ := user.Current()
	logdir := filepath.Join(usr.HomeDir, "mutter")
	c := &server.Config{logdir}
	s := server.NewServer(c)
	http.ListenAndServe("0.0.0.0:8888", s)
}

func main_old() {
	usr, _ := user.Current()
	logdir := filepath.Join(usr.HomeDir, "mutter")
	logs, err := service.GetLogFiles(logdir)
	if err != nil {
		log.Fatal(err)
	}
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
