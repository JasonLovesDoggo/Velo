package utils

import "time"

func Uint64Ptr(n uint64) *uint64 {
	return &n
}

func DurationPtr(d time.Duration) *time.Duration {
	return &d
}
