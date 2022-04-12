package utils

import (
	"strconv"
	"time"
)

func GetCronSpec(t time.Time) string {
	month := strconv.Itoa(int(t.Month()))
	day := strconv.Itoa(t.Day())
	hour := strconv.Itoa(t.Hour())
	minute := strconv.Itoa(t.Minute())
	second := strconv.Itoa(t.Second())
	spec := second + " " + minute + " " + hour + " " + day + " " + month + " *"
	return spec
}
