package utils

import (
	"strings"

	"github.com/charmbracelet/log"
)

func init() {
	log.SetReportCaller(true)
	log.SetReportTimestamp(true)
}

func StartLogger(conf AppData) {
	switch strings.ToUpper(conf.LogLevel) {
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.Info("Setting logs to", "Log Level:\t", log.GetLevel())
	log.Debug("Starting application...")
}
