package hlog

import (
	"fmt"
	"os"
)

const (
	green = "\033[32m"
	red   = "\033[31m"
	reset = "\033[0m"
)

func getStatusText(status bool) string {
	if status {
		return "[ " + green + "OK" + reset + " ]"
	} else {
		return "[ " + red + "FAILED" + reset + " ]"
	}
}

func log(msg string, status bool) {
	statusStr := getStatusText(status)
	fmt.Fprintf(os.Stderr, "%-40s %s\n", msg, statusStr)
}

func Error(err string) {
	log(err, false)
}

func Fatal(err string) {
	log(err, false)
	os.Exit(1)
}

func Ok(err string) {
	log(err, true)
}
