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
	assert.Equal(t, int64(2844), ce.GetEpochFromReferenceUnixTimestamp(ref))
}

func TestParseHourly(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Hourly,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(477735), ce.GetEpochFromReferenceUnixTimestamp(ref))
}

func TestParseDaily(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Daily,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(19905), ce.GetEpochFromReferenceUnixTimestamp(ref))
}
func TestParseWeekly(t *testing.T) {
	ref := refGlobal
	ce, err := NewCronExpression(ResetExpression{
		Type: Weekly,
	})
	assert.NoError(t, err)

	assert.Equal(t, int64(2844), ce.GetEpochFromReferenceUnixTimestamp(ref))
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
func TestUnixTimestamp(t *testing.T) {
	e := "00 6 * * 1" // every Monday at 6am
	e = "* * * * *"   // every minute
	e = "0 * * * *"   // hourly
	e = "0 0 1 * *"   // every 1st of month

	now := time.Now().UTC().Unix()

	initUnix := time.Unix(0, 0).UTC()
	first := cronexpr.MustParse(e).Next(initUnix)
	second := cronexpr.MustParse(e).Next(first)
	intervalSecs := second.Sub(first).Seconds()
	epoch := int64(math.Floor(float64((now-first.Unix())/int64(intervalSecs)))) + 1

	fmt.Println(now)
	fmt.Println("init:\t\t", initUnix)
	fmt.Println("first:\t\t", first)
	fmt.Println("second:\t\t", second)
	fmt.Println("interval:\t", intervalSecs)
	fmt.Println("epoch:\t\t", epoch)

	assert.True(t, true)
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
