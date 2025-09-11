package log

import (
	"fmt"
	"strings"
	"time"
)

var reset = "\033[0m"
var red = "\033[31m"
var green = "\033[32m"
var yellow = "\033[33m"
var cyan = "\033[36m"
var white = "\033[97m"

var (
	muted = false
)

func Mute() {
	muted = true
}

func logMessage(color string, msg any) {
	if muted {
		return
	}
	fmt.Println(white+getTime()+":"+color, prep(msg), reset)
}
func Error(msg any) {
	logMessage(red, msg)
}

func Warning(msg any) {
	logMessage(yellow, msg)
}

func Info(msg any) {
	logMessage(cyan, msg)
}

func Success(msg any) {
	logMessage(green, msg)
}

func getTime() string {
	return time.Now().Format("2006.01.02T15:04:05")
}

func prep(input any) string {
	stringInput := fmt.Sprint(input)
	segments := strings.Split(stringInput, "\n")
	for i := range segments {
		if segments[i] == "" || i == 0 {
			continue
		}
		segments[i] = strings.Repeat(" ", 21) + segments[i]
	}

	result := segments[0]
	for i := 1; i < len(segments); i++ {
		if segments[i] == "" {
			continue
		}
		result += "\n" + segments[i]
	}

	return result
}
