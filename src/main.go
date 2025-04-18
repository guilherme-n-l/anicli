package main

import (
	"anicli/opts"
	"os"
)

func main() {
	ctx, flags := opts.Ctx.GetContext(os.Args[1:])
	ctx.ParseFlags(flags)
}
