package main

import (
	"anicli/opts"
	"os"
)


func main() {
	if len(os.Args) < 2 {
		opts.Version()
		os.Exit(1)
	}

	opts.ParseFlags()
}
