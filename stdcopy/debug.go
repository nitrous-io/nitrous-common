package stdcopy

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func logf(level string, format string, a ...interface{}) {
	// Retrieve the stack infos
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "<unknown>"
		line = -1
	} else {
		file = file[strings.LastIndex(file, "/")+1:]
	}

	fmt.Fprintf(os.Stderr, fmt.Sprintf("[%s] %s:%d %s\n", level, file, line, format), a...)
}

// Debug function, if the debug flag is set, then display. Do nothing otherwise
// If Docker is in damon mode, also send the debug info on the socket
func Debugf(format string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		logf("debug", format, a...)
	}
}
