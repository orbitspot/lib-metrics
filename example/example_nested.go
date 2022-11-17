package example

import (
	e "errors"

	"github.com/orbitspot/lib-logs/pkg/errors"
)

func Level1() error {
	return Level2()
}

func Level2() error {
	err := e.New("This is my error")
	return errors.WithStack(err)
}