package support

import (
	"math"
	"math/rand"
	"time"
)

func WithExpBackoff(connector func() (conn interface{}, err error), maxTime time.Duration) (conn interface{}, err error) {
	var count uint8
	max := 1000
	min := 0

	for {
		conn, err = connector()
		if err == nil {
			break
		}
		randOffset := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
		primaryTime := time.Duration(int(math.Pow(2, float64(count)))*1000) * time.Millisecond
		var sleepTime time.Duration
		sleepTime = primaryTime + randOffset
		if sleepTime > maxTime {
			sleepTime = maxTime
		}
		time.Sleep(sleepTime)
		count++
	}
	return conn, nil
}
