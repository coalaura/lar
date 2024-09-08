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

	verbose = arguments.GetNamedAs("v", "verbose", false)

	if arguments.GetNamedAs("u", "unpack", false) {
		unpack()
	} else {
		pack()
	}
}
