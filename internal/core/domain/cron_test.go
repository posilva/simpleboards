package domain

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/gorhill/cronexpr"

	"github.com/stretchr/testify/assert"
)

// fix the reference timestamp
const refGlobal int64 = 1719848640

func TestParseCustom(t *testing.T) {
	e := "00 6 * * 1"
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type:           Custom,
		CronExpression: e,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(2847), ce.GetEpochFromReferenceUnixTimestamp(ref))
}

func TestParseHourly(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Hourly,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(477930), ce.GetEpochFromReferenceUnixTimestamp(ref))
}

func TestParseDaily(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Daily,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(19915), ce.GetEpochFromReferenceUnixTimestamp(ref))
}
func TestParseWeekly(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Weekly,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(2847), ce.GetEpochFromReferenceUnixTimestamp(ref))
}

func TestGetNexFromRefUTC(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Hourly,
	})
	assert.NoError(t, err)
	assert.Equal(t, time.Time(time.Date(2024, time.July, 1, 16, 0, 0, 0, time.UTC)),
		ce.GetNexFromRefUTC(time.Unix(ref, 0)))
}
func TestGetNexTimestampFromRefUTC(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Hourly,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1719849600),
		ce.GetNexTimestampFromRefUTC(time.Unix(ref, 0)))
}

func TestEpochByRef(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Hourly,
	})
	assert.NoError(t, err)

	epoch := ce.GetEpochFromReferenceUnixTimestamp(ref)
	assert.Equal(t, int64(477930), epoch)

}

func TestEpochByRefBeforeCached(t *testing.T) {
	ref := refGlobal
	before := time.Unix(ref, 0).Add(-60 * time.Minute).UTC().Unix()
	fmt.Println(before, ref, time.Now().UTC())
	ce, err := NewCronExpression(ResetExpression{
		Type: Hourly,
	})
	assert.NoError(t, err)

	epoch := ce.GetEpochFromReferenceUnixTimestamp(before)
	assert.Equal(t, int64(477930), epoch)

}

func BenchmarkPrimeNumbers(b *testing.B) {

	for i := 0; i < b.N; i++ {
		e := "0 0 1 * *" // every 1st of month

		now := time.Now().UTC().Unix()

		initUnix := time.Unix(0, 0).UTC()
		first := cronexpr.MustParse(e).Next(initUnix)
		second := cronexpr.MustParse(e).Next(first)
		intervalSecs := second.Sub(first).Seconds()
		_ = int64(math.Floor(float64((now-first.Unix())/int64(intervalSecs)))) + 1

	}
}
