package main

import (
	logo "github.com/sirupsen/logrus"
	"os"
)

func init() {
	logo.SetOutput(os.Stdout)
	logo.SetLevel(logo.DebugLevel)
	logo.SetReportCaller(true)
}

var log string
