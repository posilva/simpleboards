package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/gorhill/cronexpr"
)

// CronExpression data
type CronExpression struct {
	expr     *cronexpr.Expression
	first    time.Time
	second   time.Time
	interval int64
}

// NewCronExpression creates a cron expression from a reset expression
func NewCronExpression(reset ResetExpression) (CronExpression, error) {
	var e string
	switch reset.Type {
	case Hourly:
		e = "0 * * * *"
	case Daily:
		e = "0 0 * * *"
	case Weekly:
		e = "0 0 * * 7"
	case Monthly:
		e = "0 0 1 * *"
	default:
		e = reset.CronExpression
	}

	initUnix := time.Unix(0, 0).UTC()
	expr, err := cronexpr.Parse(e)
	if err != nil {
		return CronExpression{}, fmt.Errorf("failed to parse cron expression '%v': %v", e, err)
	}

	first := expr.Next(initUnix)
	second := expr.Next(first)
	intervalSecs := second.Sub(first).Seconds()

	return CronExpression{
		expr:     expr,
		first:    first,
		second:   second,
		interval: int64(intervalSecs),
	}, nil
}

// GetEpochFromReferenceUnixTimestamp calculates the epoch based on a cron expression and a ref unix timestamp
func (e *CronExpression) GetEpochFromReferenceUnixTimestamp(ref int64) int64 {
	return int64(math.Floor(float64((ref-e.first.Unix())/int64(e.interval)))) + 1
}

// GetNexFromNowUTC returns the next time after the current UTC timestamp
func (e *CronExpression) GetNexFromNowUTC() time.Time {
	return e.GetNexFromRefUTC(time.Now().UTC())
}

// GetNexFromRefUTC returns the next time after the reference timestamp
func (e *CronExpression) GetNexFromRefUTC(ref time.Time) time.Time {
	return e.expr.Next(ref.UTC())
}

// GetNexTimestampFromRefUTC returns the next time after the reference timestamp
func (e *CronExpression) GetNexTimestampFromRefUTC(ref time.Time) int64 {
	return e.expr.Next(ref.UTC()).Unix()
}
