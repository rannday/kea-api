package utils

import (
	"fmt"
	"time"
)

func HumanDuration(seconds int) string {
	d := time.Duration(seconds) * time.Second

	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	return fmt.Sprintf("%dh %dm %ds", h, m, s)
}
