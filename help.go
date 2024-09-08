package main

import (
	"os"

	"github.com/coalaura/arguments"
)

func help() {
	if !arguments.GetNamedAs("h", "help", false) {
		return
	}

	verbose = true

	info(" _")
	info("| |__ _ _ _")
	info("| / _` | '_|")
	info("|_\\__,_|_| %s", Version)

	info("\nlar [options]")

	info(" - h / help:    Show this help page")
	info(" - i / input:   Input file(s) to pack/unpack")
	info(" - o / output:  Output file/directory to pack/unpack to")
	info(" - v / verbose: Verbose mode")
	info(" - u / unpack:  Unpack instead of pack")
	info(" - t / threads: Number of threads to use (when packing)")

	info("\nExample: lar -i *.txt -o text.lar -t 4 -v")

	os.Exit(0)
}
