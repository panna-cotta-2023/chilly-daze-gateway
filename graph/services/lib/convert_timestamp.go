package lib

import (
	"log"
	"time"
)

func CovertTimestampString(timesString string) string {
	timestamp, err := time.Parse(time.RFC3339, timesString)
	if err != nil {
		log.Println("time.Parse error:", err)
		return ""
	}
	return timestamp.Format("2006-01-02T15:04:05+09:00")
}
