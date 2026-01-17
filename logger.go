package goddb

import (
	"time"

	"github.com/gookit/color"
)

var (
	colorYellow      = color.Yellow.Printf
	colorGreen       = color.Green.Printf
	colorRed         = color.Red.Printf
	colorCyan        = color.Cyan.Printf
	colorLightYellow = color.LightYellow.Printf
	colorLightBlue   = color.FgLightBlue.Printf
	colorLightGreen  = color.LightGreen.Printf
	colorLightCyan   = color.LightCyan.Printf
)

func InfoLog(log CustomLogParmas) {

	sendLog("INFO", log, colorLightBlue)
}

func DebugLog(log CustomLogParmas) {
	sendLog("debug", log, colorLightYellow)
}

func ErrorLog(log CustomLogParmas) {
	sendLog("ERROR", log, colorRed)
}

func sendLog(level string, log CustomLogParmas, printColor func(format string, a ...any)) {

	printColor("[%s] %s\n", level, log.ph)

	for k, v := range log.msg {
		printColor("\t %s: %v  ", k, v)
	}

	if log.err != nil {
		printColor("\t %s\n", log.err)
	}

	printColor("\t %s\n", time.Now().Format(time.RFC3339))
}
