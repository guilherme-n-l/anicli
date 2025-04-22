package utils

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type slug byte
type slugSet map[slug]bool

const (
	BoolSlug slug = iota
	NumberSlug
	StringSlug
)

type Context struct {
	Name           string
	Description    string
	SubCtxs        []*Context
	Slugs          slugSet
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

func NewSlugSet(slugs ...slug) *slugSet {
	ss := slugSet{}

	for _, s := range slugs {
		ss[s] = true
	}

	return &ss
}

func parseSlug(s string) slug {
	lower := strings.ToLower(s)

	if lower == "true" || lower == "false" {
		return BoolSlug
	}

	if _, err := strconv.Atoi(s); err == nil {
		return NumberSlug
	}

	return StringSlug
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

func (ss slugSet) hasSlug(s slug) bool {
	return ss[s]
}

func NewContext(name string, description string, subCtxs *[]*Context, ss *slugSet) Context {
	var fs *flag.FlagSet

	if ss == nil {
		ss = &slugSet{}
	}

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
		Slugs:          *ss,
		Fs:             fs,
		Flags:          map[*bool]func(){},
	}

	ctx.DefaultHandler = ctx.defaultHelp
	ctx.Fs.Usage = ctx.defaultHelp

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
		for _, i := range ctx.Fs.Args() {
			if !ctx.Slugs.hasSlug(parseSlug(i)) {
				ctx.invalidArgumentExit(i)
			}
		}
	}
}

func (ctx Context) AddBoolFlags(boolFlags ...BoolFlag) {
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

func (ctx Context) defaultHelp() {
	fmt.Printf("Usage of anicli %s:\n", ctx.Name)
	ctx.Fs.PrintDefaults()
	ctx.PrintSubContexts()
	os.Exit(1)
}

func (ctx Context) GetContext(args []string) (*Context, []string) {
	if len(args) == 0 {
		return &ctx, args
	}

	for _, i := range ctx.SubCtxs {
		if i.Name == args[0] {
			return i.GetContext(args[1:])
		}
	}

	return &ctx, args
}
