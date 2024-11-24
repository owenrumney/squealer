package main

import (
	"os"

	"github.com/owenrumney/squealer/internal/app/squealer/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {

	if err := cmd.Root().Execute(); err != nil {
		fail(err)
	}
}

func fail(err error) {
	log.WithError(err).Error(err.Error())
	os.Exit(-1)
}
