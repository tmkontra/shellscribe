package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/tmkontra/shellscribe/internal/server"
	"github.com/tmkontra/shellscribe/internal/shell"
)

const (
	DefaultPort = 7819
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
		var port int
		portString, ok := os.LookupEnv("SHELLSCRIBE_PORT")
		if !ok {
			port = DefaultPort
		} else {
			parsed, err := strconv.ParseInt(portString, 10, 64)
			if err != nil {
				log.Printf("Invalid port configured ('%s') using 7819 instead", portString)
				port = DefaultPort
			} else {
				port = int(parsed)
			}
		}
		c := server.NewConfig(
			cfg.LogDir(),
			port,
		)
		runServer(c)
	} else if os.Args[1] == "setup" {
		if len(os.Args) < 3 {
			log.Fatalf("must provide command string to setup")
			return
		}
		runSetup(cfg, os.Args[2])
	}
}

func runServer(cfg *server.Config) {
	s := server.NewServer(cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("server listening on: %s", addr)
	http.ListenAndServe(addr, s)
}

func runSetup(cfg *shell.Config, command string) {
	result, err := shell.SetupCommand(cfg, command)
	if err != nil {
		log.Fatalf("setup error: %s\n", err)
		return
	}
	fmt.Println(result.OutputFile())
}
