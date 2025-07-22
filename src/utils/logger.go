// --- utils/logger.go ---
package utils

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[Telemetry] ", log.LstdFlags)
