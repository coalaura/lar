package main

import (
	"github.com/coalaura/arguments"
)

var verbose bool

type File struct {
	Path  string
	Perms uint16
	Size  uint64
}

func init() {
	arguments.Parse()
}

func main() {
	help()

	verbose = arguments.Bool("v", "verbose", false)

	if arguments.Bool("u", "unpack", false) {
		unpack()
	} else {
		pack()
	}
}
