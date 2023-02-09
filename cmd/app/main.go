package main

import (
	"animal-chipization/internal/app"

	"github.com/sirupsen/logrus"
)

func main() {
	err := app.Run()
	logrus.Fatal(err)
}
