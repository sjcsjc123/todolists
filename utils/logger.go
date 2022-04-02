package utils

import (
	"TodoLists/common/config"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

var Logger = logrus.New()

func WriteLogger() {
	pathUrl := config.Conf.GetString("logInfo.path")
	info, err := rotatelogs.New(
		path.Join(pathUrl, "%Y-%m-%d info.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	errorInfo, err1 := rotatelogs.New(
		path.Join(pathUrl, "%Y-%m-%d error.log"),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil || err1 != nil {
		Logger.WithError(err).Error("unable to write logs")
		return
	}
	Logger.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.InfoLevel:  info,
			logrus.ErrorLevel: errorInfo,
			logrus.FatalLevel: info,
		}, &LogFormatter{},
	))

}

type LogFormatter struct {
}

func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	entry.Logger.SetOutput(os.Stdout)
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}
