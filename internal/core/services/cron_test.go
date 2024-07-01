package services

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/posilva/simpleboards/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	e := "00 6 * * 1"
	v := ParseCronIntervalInSeconds(e)
	assert.Equal(t, float64(604800), v)
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
	_, epoch2, _ := GetLeaderboardNameWithEpoch("", domain.ResetExpression{Type: domain.Hourly})

	fmt.Println(now)
	fmt.Println("init:\t\t", initUnix)
	fmt.Println("first:\t\t", first)
	fmt.Println("second:\t\t", second)
	fmt.Println("interval:\t", intervalSecs)
	fmt.Println("epoch:\t\t", epoch)
	fmt.Println("epoch2:\t\t", epoch2)

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
		_, _, _ = GetLeaderboardNameWithEpoch("", domain.ResetExpression{Type: domain.Hourly})
	}
}
