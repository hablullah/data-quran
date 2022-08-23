package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := (App()).Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}
