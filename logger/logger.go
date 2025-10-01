package logger

import "log"

func Error(err error) {
	log.Println("Error:", err)
}

func ErrorMsg(msg string) {
	log.Println("Error:", msg)
}

func Info(msg string) {
	log.Println("Info:", msg)
}

type Log struct {
	Disabled  bool   `json:"disabled"`
	Level     string `json:"level"`
	Output    string `json:"output"`
	Timestamp bool   `json:"timestamp"`
}

var LogConf = Log{
	Disabled:  false,
	Level:     "error",
	Output:    "",
	Timestamp: true,
}
