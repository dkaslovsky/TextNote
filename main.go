package main

import (
	"log"

	"github.com/dkaslovsky/textnote/cmd"
	"github.com/dkaslovsky/textnote/pkg/config"
)

const name = "textnote"

var version string // set by build ldflags

func main() {
	log.SetFlags(0)

	err := config.InitApp()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run(name, version)
	if err != nil {
		log.Fatal(err)
	}
}
