package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

func prettyPrint(debugFilename string, args ...interface{}) {
	var caller string
	file, fileErr := os.Create(debugFilename)
	if fileErr != nil {
		fmt.Println(fileErr)
		return
	}

	timeNow := time.Now().Format("01-02-2006 15:04:05")
	prefix := fmt.Sprintf("[%s] %s -- ", "PrettyPrint", timeNow)
	_, fileName, fileLine, ok := runtime.Caller(1)

	if ok {
		caller = fmt.Sprintf("%s:%d", fileName, fileLine)
	} else {
		caller = ""
	}

	fmt.Fprintf(file, "\n%s%s\n", prefix, caller)

	if len(args) == 2 {
		label := args[0]
		value := args[1]

		s, _ := json.MarshalIndent(value, "", "\t")
		fmt.Fprintf(file, "%s%s: %s\n", prefix, label, string(s))
	} else {
		s, _ := json.MarshalIndent(args, "", "\t")
		fmt.Fprintf(file, "%s%s\n", prefix, string(s))
	}
}
