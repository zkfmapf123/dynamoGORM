package goddb

import (
	"time"

	"github.com/gookit/color"
	dg "github.com/zkfmapf123/donggo"
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

	// msg 안의 error 타입을 문자열로 변환
	for k, v := range log.msg {
		if err, ok := v.(error); ok {
			log.msg[k] = err.Error()
		}
	}

	_, jsonStr, _ := dg.JsonStringify(log.msg)

	printColor("[%s] : %s\n", "message", jsonStr)
	printColor("\t %s\n", time.Now().Format(time.RFC3339))
}
