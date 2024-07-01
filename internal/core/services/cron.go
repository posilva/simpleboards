package services

import (
	"time"

	"github.com/gorhill/cronexpr"
)

// ParseCronIntervalInSeconds ...
func ParseCronIntervalInSeconds(expr string) float64 {
	next2 := cronexpr.MustParse(expr).NextN(time.Unix(time.Now().Unix(), 0), 2)
	secs := next2[1].Sub(next2[0]).Seconds()
	return secs
}
