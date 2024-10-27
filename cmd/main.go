package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tmkontra/shellscribe/internal/server"
	"github.com/tmkontra/shellscribe/internal/shell"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("must specify 'server' or 'setup'")
		return
	}

	cfg, err := shell.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s\n", err)
		return
	}

	if os.Args[1] == "server" {
		runServer(cfg)
	} else if os.Args[1] == "setup" {
		if len(os.Args) < 3 {
			log.Fatalf("must provide command string to setup")
			return
		}
		runSetup(cfg, os.Args[2])
	}
}

func runServer(cfg *shell.Config) {
	c := &server.Config{
		Directory: cfg.LogDir(),
	}
	s := server.NewServer(c)
	http.ListenAndServe("0.0.0.0:8888", s)
}

func runSetup(cfg *shell.Config, command string) {
	result, err := shell.SetupCommand(cfg, command)
	if err != nil {
		log.Fatalf("setup error: %s\n", err)
		return
	}
	fmt.Println(result.OutputFile())
}
