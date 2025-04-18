package main

import (
	"os"

	"anicli/opts"
)

func main() {
	ctx, flags := opts.Ctx.GetContext(os.Args[1:])
	ctx.ParseFlags(flags)
}
