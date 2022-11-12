package log

import "fmt"

// Basic log levels
var (
	Error = message(Red, '!')
	Warn  = message(Teal, ' ')
	Info  = message(Green, '*')
)

// Declare all the ANSI colors
var (
	Black   = color("\033[1;30m%s\033[0m")
	Red     = color("\033[1;31m%s\033[0m")
	Green   = color("\033[1;32m%s\033[0m")
	Yellow  = color("\033[1;33m%s\033[0m")
	Purple  = color("\033[1;34m%s\033[0m")
	Magenta = color("\033[1;35m%s\033[0m")
	Teal    = color("\033[1;36m%s\033[0m")
	White   = color("\033[1;37m%s\033[0m")
)

// color returns a function that can be used to create a function
// that can be used with Printf or Sprintf to produce colored output (ANSI).
// This is done by wrapping the format string. For example, if you want to turn
// some text green, you can do something like:
//
//	fmt.Println(log.Green("This message will be green: %s"), "woohoo!")
//
// or, more generally, if you want to write an Info, Warn, or Error message:
//
//	log.Info("this is a log message: %s", val)
//	log.Warn("this is a warning message: %s", val)
//	log.Error("this is an error message: %s", val)
//
// These functions are defined at the top of the file in variable declarations
func color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

// message is used to create various ANSI colors. See the variable declarations
// at the top of the file.
func message(colorFunc func(...interface{}) string, symbol rune) func(string, ...interface{}) {
	return func(message string, args ...interface{}) {
		fmt.Printf(colorFunc("[%c] %s\n"), symbol, fmt.Sprintf(message, args...))
	}
}
