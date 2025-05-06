package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogger() {

	// Проверяем, существует ли папка logs, если нет — создаем
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		Log.Warn("Не удалось создать папку для логов, логируем в stdout")
	}

	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.SetOutput(file)
	} else {
		Log.SetOutput(os.Stdout)
		Log.Warn("Не удалось открыть файл лога, логируем в stdout")
	}

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.SetLevel(logrus.InfoLevel)
}
