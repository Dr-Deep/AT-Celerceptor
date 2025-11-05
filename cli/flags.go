package main

import (
	"fmt"

	celerceptor "github.com/Dr-Deep/AT-Celerceptor"
)

func flagUsage() {
	fmt.Printf(
		`Usage:
	%s <action> [args]

Actions:
	%s run		:	selbstständig nachbuchen
	%s get		:	Aktuelles Daten-Volumen erfragen
	%s nachbuchen	:	Daten-Volumen nachbuchen
	%s version		:	version info
`, prog, prog, prog, prog, prog)
}

/*
var (

	// Flags
	logFilePath = flag.String(
		"logfile",
		"",
		"file to redirect logs to",
	)
	logLevel = flag.String(
		"loglevel",
		"info",
		"log level ('debug', 'info', 'error', 'fatal', 'none')",
	)
	configFilePath = flag.String(
		"config-file",
		"./config.yml",
		"configuration file",
	)
)
*/

// normal: loop: datavol? & nachbuchen
func flagModeRun() {
	var (
		logger = initLogger(defaultLogFile, defaultLogLevel)
		cfg    = initConfig(defaultCfgFile)
	)
	defer logger.Close()

	var c = celerceptor.NewCelerceptor(logger, cfg)
	if err := c.Launch(); err != nil {
		panic(err)
	}
}

// nur datenvolumen ausgeben
func flagModeGetDataVol() {}

// nur nachbuchen
func flagModeKaufDataVol() {}
