package main

import (
	"fmt"
	"os"

	"github.com/smoothie-go/smgo-manager/commands"
)

func main() {
	if len(os.Args)-1 < 1 {
		fmt.Printf(`Usage:
%s set`, os.Args[0])
		return
	}
	switch os.Args[1] {
	case "set":
		commands.Set()
	default:
		fmt.Printf(`Usage:
%s set`, os.Args[0])
	}
}
