package main

import (
	"fmt"
	"os"
)

func info(fm string, args ...interface{}) {
	if !verbose {
		return
	}

	fmt.Printf(fm+"\n", args...)
}

func fatalf(code int, fm string, args ...interface{}) {
	fmt.Printf("ERROR: "+fm+"\n", args...)

	os.Exit(code)
}
