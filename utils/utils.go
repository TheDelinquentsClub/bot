package utils

import (
	"fmt"
	"time"
)

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func GetDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%02d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}
