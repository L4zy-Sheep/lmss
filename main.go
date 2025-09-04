package main

import (
	"lmss/cmd"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout) //test
	log.SetPrefix("[lmss] ")
	cmd.RootCmd.Execute()
}
