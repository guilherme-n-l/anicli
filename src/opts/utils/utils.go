package utils

import "flag"

type FlagContext struct {
	Fs    *flag.FlagSet
	Flags map[*bool]func()
}

func NewBoolFlag(longhand string, shorthand string, value bool, desc string, fs *flag.FlagSet) *bool {
	var newFlag *bool

	if fs == nil {
		newFlag = flag.Bool(longhand, value, desc)
		if shorthand != "" {
			flag.BoolVar(newFlag, shorthand, value, desc)
		}
	} else {
		newFlag = (*fs).Bool(longhand, value, desc)
		if shorthand != "" {
			(*fs).BoolVar(newFlag, shorthand, value, desc)
		}
	}

	return newFlag
}
