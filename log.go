package crescent

import (
	"github.com/Sirupsen/logrus"
)

var log = logrus.New()

func Logger() *logrus.Logger {
	return log
}
