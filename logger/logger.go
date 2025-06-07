package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("|%s| |service-system| |%s|  :  %s", timestamp, entry.Level, entry.Message))
	if len(entry.Data) > 0 {
		b.WriteString(" | ")
		for key, value := range entry.Data {
			b.WriteString(fmt.Sprintf("%s=%v ", key, value))
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func InitLogger() {
	Log = logrus.New()

	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     7,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	Log.SetOutput(multiWriter)

	Log.SetFormatter(new(CustomFormatter))
	Log.SetLevel(logrus.InfoLevel)
}
