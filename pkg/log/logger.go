package log

import (
	"fmt"
	"github.com/fatih/color"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	maxFileNameSize = 12
	LDebug          = "DEBUG"
	LInfo           = "INFO "
	LWarn           = "WARN "
	LError          = "ERROR"
	LFatal          = "FATAL"
	LTrace          = "TRACE"
	LPanic          = "PANIC"
)

var (
	internalLog *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmsgprefix)
	logLevelApp             = os.Getenv("LOG_LEVEL")
	logCodeName             = os.Getenv("LOG_CODE_NAME")
	environment             = os.Getenv("ENVIRONMENT")
	appName                 = os.Getenv("APP_NAME")
	logCode                 = true
)

// Init Logging
func Init() {
	internalLog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmsgprefix)
	if logLevelApp == "" {
		if environment == "DEV" || environment == "DEVELOPMENT" {
			logLevelApp = LDebug
		} else {
			logLevelApp = LInfo
		}
	}
	if logCodeName == "FALSE" {
		logCode = false
	}
	// Tries to start Sentry for error observability
	startSentry()
	Info("Logging STARTED! [app: %s]", appName)
}

// Trace Specific Log for Traces
func Trace(msg string, args ...interface{}) {
	if logLevelApp == LDebug || logLevelApp == LInfo {
		logger(LTrace, fmt.Sprintf(msg, args...))
	}
}

// Debug level for Logging
func Debug(msg string, args ...interface{}) {
	if logLevelApp == LDebug {
		logger(LDebug, fmt.Sprintf(msg, args...))
	}
}

// Info level for Logging
func Info(msg string, args ...interface{}) {
	if logLevelApp == LDebug || logLevelApp == LInfo {
		logger(LInfo, fmt.Sprintf(msg, args...))
	}
}

// Warn level for Logging
func Warn(msg string, args ...interface{}) {
	if logLevelApp == LDebug || logLevelApp == LInfo || logLevelApp == LWarn {
		logger(LWarn, fmt.Sprintf(msg, args...))
	}
}

// Error Default Logging based on an error
func Error(err error) {
	newErr := errors.WithStack(err).(*errors.Error)
	message := fmt.Sprintf("%s", newErr)
	logger(LError, message)
	loggerStackTrace(LError, string(newErr.Stack()))
	logErrorSentry(err)
}

// Errorf Logging based on an error and a custom message for additional context and variables
func Errorf(err error, msg string, args ...interface{}) {
	newErr := errors.WithStack(err).(*errors.Error)
	message := fmt.Sprintf("%s - Error: %s", fmt.Sprintf(msg, args...), newErr)
	logger(LError, message)
	loggerStackTrace(LError, string(newErr.Stack()))
	logErrorSentry(err)
}

// ErrorNew Creates a new Error, log it and return the generated error
func ErrorNew(msg string, args ...interface{}) error {
	err := errors.New(fmt.Sprintf(msg, args...))
	newErr := errors.WithStack(err).(*errors.Error)
	message := fmt.Sprintf("%s", newErr)
	logger(LError, message)
	loggerStackTrace(LError, string(newErr.Stack()))
	logErrorSentry(err)
	return err
}

// ErrorHandler Just log an error and return the same error
func ErrorHandler(err error) error {
	if err != nil {
		newErr := errors.WithStack(err).(*errors.Error)
		message := fmt.Sprintf("%s", newErr)
		logger(LError, message)
		loggerStackTrace(LError, string(newErr.Stack()))
		logErrorSentry(err)
		return err
	}
	return nil
}

// Fatal Error Default Logging based on an error
// It should be used when...
//
//	an error happens in any func Init()
//	an error happens of which is irrecoverable
//	an error occurs during a process which might not be reversible
func Fatal(err error, args ...bool) {
	newErr := errors.WithStack(err).(*errors.Error)
	message := fmt.Sprintf("%s", newErr)
	logger(LFatal, message)
	loggerStackTrace(LFatal, string(newErr.Stack()))
	logErrorSentry(err)
	panic(err)
}

// Fatalf Logging based on a fatal error and a custom message for additional context and variables
//
//	an error happens in any func Init()
//	an error happens of which is irrecoverable
//	an error occurs during a process which might not be reversible
//
// It should be used when...
func Fatalf(err error, msg string, args ...interface{}) {
	newErr := errors.WithStack(err).(*errors.Error)
	message := fmt.Sprintf("%s - Error: %s", fmt.Sprintf(msg, args...), newErr)
	logger(LFatal, message)
	loggerStackTrace(LFatal, string(newErr.Stack()))
	logErrorSentry(err)
	panic(err)
}

// Recovery sentry util in gin
func Recovery() {
	if r := recover(); r != nil {
		Simple(LPanic, "Call to Recovery() [%s]", r)
	}
}

// Simple visualization
func Simple(level string, msg string, args ...interface{}) {
	color.Set(getColor(level))
	internalLog.Println(level + ": " + fmt.Sprintf(msg, args...))
	color.Unset()
}

// Boxed visualization with no timestamp info
func Boxed(level string, msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	size := len(message)
	Space(level)
	color.Set(getColor(level))
	internalLog.Println(level + ": " + strings.Repeat("*", size+8))
	internalLog.Println(level + ": " + "**  " + message + "  **")
	internalLog.Println(level + ": " + strings.Repeat("*", size+8))
	color.Unset()
	Space(level)
}

// Lines Helper for better logs visualization
func Lines(level string) {
	color.Set(getColor(level))
	internalLog.Println(level + ": " + strings.Repeat("*", 80))
	color.Unset()
}

// Breakline Helper for better logs visualization
func Space(level string, lines ...int) {
	color.Set(getColor(level))
	if len(lines) > 0 {
		for i := 0; i < lines[0]; i++ {
			internalLog.Println(level + ": ")
		}
	} else {
		internalLog.Println(level + ": ")
	}
	color.Unset()
}

// Just log an error and return the same error
func errorNoSentry(err error, msg string, args ...interface{}) {
	message := fmt.Sprintf("%s - Error: %s", fmt.Sprintf(msg, args...), err)
	logger(LError, message)
}

func logger(level string, message string) {
	color.Set(getColor(level))
	codeFileName := ""
	if logCode {
		codeFileName = " >> " + getCodeFileName(2)
	}
	messages := strings.Split(message, "\n")
	lines := len(messages)
	if lines > 1 {
		internalLog.Println(level + ": " + messages[0]) // First Line
		for i := 1; i < lines; i++ {
			if i == lines-1 {
				internalLog.Println(level + ":  " + messages[i] + codeFileName) // Last Line
			} else {
				internalLog.Println(level + ":  " + messages[i]) // In between Lines
			}
		}
	} else {
		internalLog.Println(level + ": " + message + codeFileName)
	}
	//internalLog.Println(fmt.Sprintf("%s %s: "+message, getCodeFileName(2), level))
	color.Unset()
}

func loggerStackTrace(level string, message string) {
	color.Set(getColor(level))
	internalLog.Println(level + ": StackTrace:")
	messages := strings.Split(message, "\n")
	for i := 1; i < len(messages); i++ {
		internalLog.Println(level + ":   " + messages[i])
	}
	color.Unset()
}

func getColor(level string) color.Attribute {
	switch level {
	case LInfo:
		return color.FgHiBlue
	case LDebug:
		return color.FgGreen
	case LWarn:
		return color.FgYellow
	case LTrace:
		return color.FgHiMagenta
	case LError:
		return color.FgHiRed
	case LFatal:
		return color.FgRed
	default:
		return color.Reset
	}
}

func getCodeFileName(skip int) string {
	// Pay attention for refactorings -> we are in 3rd level of function caller
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		file = "???"
		line = 0
	} else {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
	}
	file = fmt.Sprintf("%s:%d", file, line)
	return file
}
