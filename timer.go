package main

import "time"

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

func (t *Timer) Run(onTickFn, onCompleteFn func() error) {
	for {
		select {
		case <-t.ticker.C:
			t.remaining -= time.Second
			if err := onTickFn(); err != nil {
				panic(err)
			}
		case <-t.timer.C:
			if err := onCompleteFn(); err != nil {
				panic(err)
			}
		}
	}
}
