package telemetry

import "log"

func LogInfo(msg string) {
	log.Printf("[INFO] %s", msg)
}

func LogError(msg string) {
	log.Printf("[ERROR] %s", msg)
}

func LogDebug(msg string) {
	log.Printf("[DEBUG] %s", msg)
}