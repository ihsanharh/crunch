package utils

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandomBool() bool {
	return rand.Int63()&1 == 0
}

func FormatTime(duration time.Duration) string {
	totalSeconds := int64(duration.Seconds())
	days := totalSeconds / 86400
	hours := totalSeconds % 86400 / 3600
	minutes := totalSeconds % 3600 / 60
	seconds := totalSeconds % 60

	if days > 0 {
		return Fmt("%02d:%02d:%02d:%02d", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return Fmt("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return Fmt("%02d:%02d", minutes, seconds)
}

func ParseDuration(duration string) (time.Duration, error) {
	if duration, err := time.ParseDuration(duration); err == nil {
		return duration, nil
	}

	if _, err := strconv.Atoi(strings.ReplaceAll(duration, ":", "")); err != nil || len(duration) > 8 {
		return 0, errors.New("invalid duration")
	}

	parts := strings.Split(duration, ":")
	result := Fmt("%ss", parts[0])

	if len(parts) == 3 {
		result = Fmt("%sh%sm%ss", parts[0], parts[1], parts[2])
	}
	if len(parts) == 2 {
		result = Fmt("%sm%ss", parts[0], parts[1])
	}

	return time.ParseDuration(result)
}

func StringArrayContains(array []string, term string) bool {
	for _, item := range array {
		if item == term {
			return true
		}
	}

	return false
}

func IntegerArrayContains(array []int, number int) bool {
	for _, i := range array {
		if i == number {
			return true
		}
	}
	return false
}