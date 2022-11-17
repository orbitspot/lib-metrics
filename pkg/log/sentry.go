package log

import (
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"os"
	_ "github.com/joho/godotenv/autoload"
)

var (
	// sentryDsn sentry DNS
	sentryDsn = os.Getenv("SENTRY_DSN")
)

// startSentry initialize Sentry
func startSentry() {
	if sentryDsn != "" {
		if err := raven.SetDSN(sentryDsn); err != nil {
			errorNoSentry(err, "Sentry WAS NOT STARTED because exists a error in Sentry or in config [DSN: %s]", sentryDsn)
		} else {
			Info("Sentry STARTED for [DSN: %s]", sentryDsn)
		}
	} else {
		Warn("Sentry WAS NOT STARTED because there is no config")
	}
}

// LogError log an error
func logErrorSentry(err error) {
	if sentryDsn != "" {
		raven.CaptureError(err, nil)
	}
}

// capturePanic capture panic from Sentry
func capturePanic(f func()) {
	if sentryDsn != "" {
		raven.CapturePanicAndWait(f, nil)
	}
}

// recovery recovery Sentry
func recovery() {
	if sentryDsn != "" {
		if r := recover(); r != nil {
			jsonError, _ := json.Marshal(r)
			err := errors.New(string(jsonError))
			logErrorSentry(err)
			errorNoSentry(err, "Sentry WAS NOT RECOVERED because exists a error in Sentry or in config [DSN: %s]", sentryDsn)
		}
	}
}