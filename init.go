package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	logFile *os.File
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		cobra.CheckErr(err)
	}
	// open a file
	logFile, err = os.OpenFile(filepath.Join(home, "typora-pic-upload.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}

	// don't forget to close it
	// defer f.Close()

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	logrus.SetOutput(logFile)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)
}
