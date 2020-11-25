package clock

import (
	"time"
)

type Ticks = <-chan time.Time
type Clock = *clock

type clock struct {
	ticks Ticks
}

func New(period time.Duration) Clock {
	return &clock{ticks:time.Tick(period)}
}

func (c *clock) Ticks() Ticks {
	return c.ticks
}