package main

import (
	"worker-service/cmd/cli"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	cli.Execute()
}
