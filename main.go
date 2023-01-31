package main

import (
	"fmt"
	"net-cat/cmd"
	"os"
)

func main() {
	PORT := "8080"
	switch true {
	case len(os.Args) > 2:
		fmt.Println("[USAGE]: go run . $port")
		return
	case len(os.Args) == 2:
		PORT = os.Args[1]
	}
	cmd.NetRun(PORT)
}
