package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tmkontra/shellscribe/internal/shell"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("must provide command string")
	}
	result, err := shell.SetupCommand(os.Args[1])
	if err != nil {
		log.Fatalf("setup error: %s\n", err)
	}
	fmt.Println(result.OutputFile())
}
