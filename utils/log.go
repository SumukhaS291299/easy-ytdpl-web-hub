package utils

import "github.com/charmbracelet/log"

func init() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(true)
}

func StartLogger() {
	log.Debug("Starting application...")
}
