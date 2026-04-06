package logging

import "fmt"

const (
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[1;33m"
	nc     = "\033[0m"
)

func Info(msg string)  { fmt.Printf("%sINFO%s: %s\n", yellow, nc, msg) }
func Step(msg string)  { fmt.Printf("\n%s==>%s %s\n", green, nc, msg) }
func Error(msg string) { fmt.Printf("%sERROR%s: %s\n", red, nc, msg) }
