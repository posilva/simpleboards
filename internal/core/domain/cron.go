package domain

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

type epochCache struct {
	epoch     int64
	timestamp time.Time
}

// CronExpression data
type CronExpression struct {
	expr            *cronexpr.Expression
	first           time.Time
	second          time.Time
	interval        int64
	lastEpochCached epochCache
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

	expr, err := cronexpr.Parse(e)
	if err != nil {
		return CronExpression{}, fmt.Errorf("failed to parse cron expression '%v': %v", e, err)
	}
	initUnix := time.Unix(0, 0).UTC()
	first := expr.Next(initUnix)
	second := expr.Next(first)
	intervalSecs := second.Sub(first).Seconds()

	lec := calculateCurrentEpochSince(expr, initUnix)

	ce := CronExpression{
		expr:            expr,
		first:           first,
		second:          second,
		interval:        int64(intervalSecs),
		lastEpochCached: lec,
	}
	return ce, nil
}

func calculateCurrentEpochSince(expr *cronexpr.Expression, since time.Time) epochCache {
	now := time.Now().UTC()

	next := since
	previous := next
	var epoch int64 = 1
	for next.Unix() < now.Unix() {
		previous = next
		next = expr.Next(next)
		epoch++
	}
	lec := epochCache{
		epoch:     epoch,
		timestamp: previous,
	}
	return lec
}

// GetCurrentEpoch returns the epoch based in the current Timestamp
func (e *CronExpression) GetCurrentEpoch() int64 {
	now := time.Now().UTC().Unix()

	return e.GetEpochFromReferenceUnixTimestamp(now)
}

// GetEpochFromReferenceUnixTimestamp calculates the epoch based on a cron expression and a ref unix timestamp
func (e *CronExpression) GetEpochFromReferenceUnixTimestamp(ref int64) int64 {
	epoch := e.lastEpochCached.epoch
	since := e.lastEpochCached.timestamp

	fmt.Println(epoch, since)

	if ref >= e.lastEpochCached.timestamp.Unix() {
		ec := calculateCurrentEpochSince(e.expr, since)
		epoch = epoch + ec.epoch
		since = ec.timestamp
		e.lastEpochCached = epochCache{epoch: epoch, timestamp: since}
	} else {
		initUnix := time.Unix(0, 0).UTC()
		ec := calculateCurrentEpochSince(e.expr, initUnix)
		epoch = ec.epoch
		since = ec.timestamp
	}

	return epoch
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
