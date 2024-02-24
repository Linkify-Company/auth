package logger

import (
	"auth/pkg/errify"
	"fmt"
	"os"
	"time"
)

func saveErrify(err errify.IError, level string) {
	toFile(fmt.Sprintf("time=%s level=%s location=%s message=%q code=%d details=%q read=%q",
		time.Now().Format(time.DateTime), level, err.Location(), err.Message(), err.Code(), err.Details(), err.Read()))
}

func save(err string, level string) {
	toFile(fmt.Sprintf("time=%s level=%s error=%s", time.Now().Format(time.DateTime), level, err))
}

func toFile(err string) {
	file, e := os.OpenFile("./pkg/logs/app.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if e != nil {
		fmt.Printf("Не удалось открыть файл лога: %v", e)
		return
	}
	defer file.Close()
	_, e = file.WriteString(fmt.Sprintln(err))
	if e != nil {
		fmt.Printf("Не удалось записать логи в файл: %v", e)
		return
	}
}
