package main

import (
	"time"
)

type Timer struct {
	duration, remaining time.Duration
	ticker              *time.Ticker
	timer               *time.Timer
}

func NewTimer(duration time.Duration) Timer {
	return Timer{
		duration:  duration,
		remaining: duration,
		ticker:    time.NewTicker(time.Second),
		timer:     time.NewTimer(duration),
	}
}

func (t *Timer) OnTick(tickFn func() error) {
	for {
		select {
		case <-t.ticker.C:
			t.remaining -= time.Second
			if err := tickFn(); err != nil {
				panic(err)
			}
		case <-t.timer.C:
			return
		}
	}
}
