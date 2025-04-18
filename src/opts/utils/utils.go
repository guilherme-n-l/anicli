package utils

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Context struct {
	Name           string
	Description    string
	SubCtxs        []*Context
	Fs             *flag.FlagSet
	Flags          map[*bool]func()
	DefaultHandler func()
}

type Flag struct {
	LongHand    string
	ShortHand   string
	Description string
}

type BoolFlag struct {
	Flag
	Value   bool
	Handler func()
}

func NewBoolFlag(longhand string, shorthand string, value bool, description string, handler func()) BoolFlag {
	return BoolFlag{
		Flag: Flag{
			LongHand:    longhand,
			ShortHand:   shorthand,
			Description: description,
		},
		Handler: handler,
	}
}

func addBoolFlag(bf *BoolFlag, fs *flag.FlagSet) *bool {
	var newFlag *bool

	newFlag = (*fs).Bool(bf.LongHand, bf.Value, bf.Description)
	if bf.ShortHand != "" {
		(*fs).BoolVar(newFlag, bf.ShortHand, bf.Value, bf.Description)
	}

	return newFlag
}

func NewContext(name string, description string, subCtxs *[]*Context) Context {
	var fs *flag.FlagSet

	if subCtxs == nil {
		subCtxs = &[]*Context{}
	}

	if len(name) == 0 {
		name = "main"
		fs = flag.CommandLine
	} else {
		fs = flag.NewFlagSet(name, flag.ExitOnError)
	}

	var ctx = Context{
		Name:           name,
		Description:    description,
		SubCtxs:        *subCtxs,
		Fs:             fs,
		Flags:          map[*bool]func(){},
		DefaultHandler: fs.PrintDefaults,
	}
	return ctx
}

func (ctx Context) invalidArgumentExit(arg string) {
	const InvalidArgumentError = "Invalid argument provided not allowed"
	log.Printf("%s: %s\n", InvalidArgumentError, arg)
	ctx.Fs.PrintDefaults()
	ctx.PrintSubContexts()
	os.Exit(1)
}

func (ctx Context) preventInvalidArgs() {
	if ctx.Fs.NArg() > 0 {
		ctx.invalidArgumentExit(ctx.Fs.Arg(0))
	}
}

func isFlag(s string) bool {
	if len(s) == 0 {
		return false
	}

	if len(s) > 1 && s[:2] == "--" || s[:1] == "-" {
		return true
	}

	return false
}

func (ctx Context) AddBoolFlags(boolFlags []BoolFlag) {
	for _, i := range boolFlags {
		ctx.Flags[addBoolFlag(&i, ctx.Fs)] = i.Handler
	}
}

func (ctx Context) ParseFlags(args []string) {
	ctx.Fs.Parse(args)

	ctx.preventInvalidArgs()

	for f, handler := range ctx.Flags {
		if *f {
			handler()
		}
	}

	ctx.DefaultHandler()
}

func (ctx Context) PrintSubContexts() {
	for _, i := range ctx.SubCtxs {
		fmt.Printf("  %s\n        %s\n", i.Name, i.Description)
	}
}

func (ctx Context) GetContext(args []string) (*Context, []string) {
	if len(args) == 0 || isFlag(args[0]) {
		return &ctx, args
	}

	for _, i := range ctx.SubCtxs {
		if i.Name == args[0] {
			return i.GetContext(args[1:])
		}
	}

	ctx.invalidArgumentExit(args[0])

	return nil, nil
}
