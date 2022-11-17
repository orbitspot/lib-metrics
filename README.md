# Logs GO 

Biblioteca para padronizar os LOGs em nossos projetos GO.

Tambem suporta os logs de ERROS no SENTRY.

## Como utilizar essa biblioteca em seu projeto

No diretório do seu projeto em GO, obtendoa versão mais recente disponível. Exemplo:
```
go get github.com/orbitspot/lib-logs
```

1 - Dentro de seu arquivo `.env`, adicione as variáveis de ambiente abaixo, com a URL do sentry (configurada por projeto). Exemplo:
```
APP_NAME=mytest
ENVIRONMENT=DEV
LOG_LEVEL=DEBUG
LOG_CODE_NAME=TRUE # Print Code file and line in all log lines, eg:  'main.go:34'
SENTRY_DSN=  # Observability for Errors. Format => https://XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX@sentry.domain/NN
```
**Atencão: Solicitar para o time de DEVOPS preparar o Sentry em producão e configurar o pipeline para o seu projeto**

2 - Segue um exemplo de como utilizar os logs. Lembre-se de importar nos códigos corretamente o pacote **"github.com/orbitspot/lib-logs/pkg/log"**

```
package main

import (
	"github.com/orbitspot/lib-logs/example"
	"github.com/orbitspot/lib-logs/pkg/errors"
	"github.com/orbitspot/lib-logs/pkg/log"
	_ "github.com/joho/godotenv/autoload"
)

func main(){
	// Init Logging & Sentry => check file `.env`
	log.Init()

	// USE WITH CAUTION! Just an example of an function to recover()
	defer log.Recovery()

	log.Boxed(log.LInfo,"DEFAULT Logs")
	log.Trace("Testing [var1: %s, var2: %v]", "my variable", 123)
	log.Debug("Testing [var1: %s, var2: %v]", "my variable", 123)
	log.Info("Testing [var1: %s, var2: %v]", "my variable", 123)
	log.Warn("Testing [var1: %s, var2: %v]", "my variable", 123)

	log.Boxed(log.LInfo,"MULTIPLE lines Logs")
	log.Info("Testing [var1: %s, var2: %v]\nLog line 2\nLog Line 3\nLog Last Line", "my variable", 123)

	log.Boxed(log.LInfo,"ERROR with NO additional formatted text and WITH StackTrace")
	err := errors.New("This is my error")
	log.Error(err)

	log.Boxed(log.LInfo,"ERROR with additional formatted text and WITH StackTrace")
	log.Errorf(err,"Testing Error [var1: %s, var2: %v]", "my variable", 123)

	log.Boxed(log.LInfo,"Creates a NEW ERROR, Log & Return WITH StackTrace")
	_ = log.ErrorNew("Testing my Error [var1: %s, var2: %v]", "my variable", 123)

	log.Boxed(log.LInfo,"HANDLE Log & Return ERROR WITH StackTrace")
	err = log.ErrorHandler(err)


	log.Boxed(log.LInfo,"Nested Errors")
	err = example.Level1()
	log.Error(err)

	log.Boxed(log.LInfo,"HELPERS for Logs. You must always inform Log Level.")
	log.Lines(log.LInfo) // Print a full line with asterisks and Info Color
	log.Lines(log.LInfo)
	log.Simple(log.LInfo,"Not Recommended! Simple message with no code filename and line number")
	log.Lines(log.LInfo)
	log.Lines(log.LInfo)
	log.Space(log.LInfo)   // Print a single Break Lines with Info Color
	log.Space(log.LInfo,3) // Print 3 Break lines with Info Color

	log.Boxed(log.LInfo,"FATAL with additional formatted text and WITH StackTrace")
	log.Fatalf(err,"Testing PANIC [var1: %s, var2: %v]", "my variable", 123)

	log.Info("Tutorial Finished!")
}
```
