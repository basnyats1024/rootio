package rootio

import (
	"fmt"
	"os"
)

var g_rootio_debug = false

func myprintf(format string, args ...interface{}) (n int, err error) {
	if g_rootio_debug {
		return fmt.Printf(format, args...)
	}
	return
}

func init() {
	switch os.Getenv("ROOTIO_DEBUG") {
	case "0", "":
		g_rootio_debug = false
	default:
		g_rootio_debug = true
	}
}
