package main

import (
	"flag"

	"github.com/tevino/log"
)

type opt struct {
	logLevel    log.Level
	isDebug     bool
	needExample bool
	needClean   bool
}

func (o *opt) Parse() *opt {
	var logLevelStr string
	flag.StringVar(&logLevelStr, "log", "INFO", "Log level [DEBUG|INFO|WARNING|FATAL]")
	flag.BoolVar(&o.isDebug, "debug", false, "debug mode, alias of -log DEBUG")
	flag.BoolVar(&o.needExample, "example", false, "Create directory structure with examples in the current directory")
	flag.BoolVar(&o.needClean, "clean", false, "Remove all files in the peers folder before generating")
	flag.Parse()

	o.logLevel = log.LevelFromString(logLevelStr)
	if o.isDebug {
		o.logLevel = log.DEBUG
	}
	return o
}
