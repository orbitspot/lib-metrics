package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/orbitspot/lib-metrics/example"
	"github.com/orbitspot/lib-metrics/pkg/errors"
	"github.com/orbitspot/lib-metrics/pkg/log"
)

func main() {
	// Init Logging & Sentry => check file `.env`
	log.Init()

	// USE WITH CAUTION! Just an example of an function to recover()
	defer log.Recovery()

	log.Boxed(log.LInfo, "DEFAULT Logs")
	log.Trace("Testing [var1: %s, var2: %v]", "my variable", 123)
	log.Debug("Testing [var1: %s, var2: %v]", "my variable", 123)
	log.Info("Testing [var1: %s, var2: %v]", "my variable", 123)
	log.Warn("Testing [var1: %s, var2: %v]", "my variable", 123)

	log.Boxed(log.LInfo, "MULTIPLE lines Logs")
	log.Info("Testing [var1: %s, var2: %v]\nLog line 2\nLog Line 3\nLog Last Line", "my variable", 123)

	log.Boxed(log.LInfo, "ERROR with NO additional formatted text and WITH StackTrace")
	err := errors.New("This is my error")
	log.Error(err)

	log.Boxed(log.LInfo, "ERROR with additional formatted text and WITH StackTrace")
	log.Errorf(err, "Testing Error [var1: %s, var2: %v]", "my variable", 123)

	log.Boxed(log.LInfo, "Creates a NEW ERROR, Log & Return WITH StackTrace")
	_ = log.ErrorNew("Testing my Error [var1: %s, var2: %v]", "my variable", 123)

	log.Boxed(log.LInfo, "HANDLE Log & Return ERROR WITH StackTrace")
	err = log.ErrorHandler(err)

	log.Boxed(log.LInfo, "Nested Errors")
	err = example.Level1()
	log.Error(err)

	log.Boxed(log.LInfo, "HELPERS for Logs. You must always inform Log Level.")
	log.Lines(log.LInfo) // Print a full line with asterisks and Info Color
	log.Lines(log.LInfo)
	log.Simple(log.LInfo, "Not Recommended! Simple message with no code filename and line number")
	log.Lines(log.LInfo)
	log.Lines(log.LInfo)
	log.Space(log.LInfo)    // Print a single Break Lines with Info Color
	log.Space(log.LInfo, 3) // Print 3 Break lines with Info Color

	log.Boxed(log.LInfo, "FATAL with additional formatted text and WITH StackTrace")
	log.Fatalf(err, "Testing PANIC [var1: %s, var2: %v]", "my variable", 123)

	log.Info("Tutorial Finished!")

}
