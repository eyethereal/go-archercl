package config

import (
	"github.com/op/go-logging"
)

type delayedMessage struct {
	level   logging.Level
	message string
}

var delayed = make([]*delayedMessage, 0, 5)

func logDelayed(level logging.Level, msg string) {
	dm := &delayedMessage{
		level:   level,
		message: msg,
	}
	delayed = append(delayed, dm)
}

func outputDelayedLog(lgr *logging.Logger) {

	for _, d := range delayed {
		switch d.level {
		case logging.CRITICAL:
			lgr.Critical("%s", d.message)

		case logging.ERROR:
			lgr.Error("%s", d.message)

		case logging.WARNING:
			lgr.Warning("%s", d.message)

		case logging.NOTICE:
			lgr.Notice("%s", d.message)

		case logging.INFO:
			lgr.Info("%s", d.message)

		default:
			lgr.Debug("%s", d.message)

		}
	}

	// Clear them
	delayed = make([]*delayedMessage, 0, 5)
}
