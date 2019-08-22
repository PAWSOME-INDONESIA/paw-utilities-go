package logging

import (
	"fmt"
	"io"
	"os"

	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func Log() {
	var filename = "storages/logs/logfile - " + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	Formatter := &prefixed.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}

	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		mw := io.MultiWriter(os.Stdout, f)
		logrus.SetOutput(mw)
		// log.SetOutput(f)
	}

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(Formatter)

}
