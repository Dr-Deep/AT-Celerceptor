/*
	!!! PROOF OF CONCEPT - HAFTUNGSAUSSCHLUSS !!!

Dieser Code dient rein der Illustration.
Nutzung auf eigene Gefahr. Keine Garantie für Funktionalität oder Richtigkeit.
Der Entwickler übernimmt keine Haftung für Schäden oder Folgenutzung.
Vor produktivem Einsatz rechtliche und technische Risiken prüfen.
*/

package main

import (
	"flag"
	"os"

	"github.com/Dr-Deep/AT-Celerceptor/internal"
	"github.com/Dr-Deep/AT-Celerceptor/internal/config"
	"github.com/Dr-Deep/logging-go"
)

const rwForOwnerOnlyPerm = 0o600

var (
	logger *logging.Logger
	cfg    *config.Configuration

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

func setupLogger() {
	if *logFilePath != "" {
		// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
		logFile, err := os.OpenFile(
			*logFilePath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			rwForOwnerOnlyPerm,
		)
		if err != nil {
			panic(err)
		}

		logger = logging.NewLogger(logFile)
	} else {
		logger = logging.NewLogger(os.Stdout)
	}

	switch *logLevel {
	case "debug":
		logger.Level = logging.LogDebug

	case "info":
		logger.Level = logging.LogInfo

	case "error":
		logger.Level = logging.LogError

	case "fatal":
		logger.Level = logging.LogFatal

	case "none":
		logger.Level = logging.Level(0)
	}
}

func setupConfig() {
	// #nosec G304 -- Zugriff nur auf bekannte Log- und Config-Dateien
	cfgFile, err := os.OpenFile(
		*configFilePath,
		os.O_RDONLY,
		rwForOwnerOnlyPerm,
	)
	if err != nil {
		panic(err)
	}

	cfg, err = config.UnmarshalConfigFile(cfgFile)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()
	setupLogger()
	setupConfig()
	defer logger.Close()

	// Bot
	var c = internal.NewCelerceptor(logger, cfg)
	if err := c.Launch(); err != nil {
		panic(err)
	}

	logger.Close()
	os.Exit(0)
}
