package main

import (
	"time"
)

type TimerWindow struct {
	name   string
	timer  *Timer
	flower *Flower
}

func NewTimerWindow(name string, duration time.Duration) TimerWindow {
	timer := NewTimer(duration)
	flower := NewFlower("hyacinth")

	return TimerWindow{
		name:   name,
		timer:  &timer,
		flower: &flower,
	}
}
