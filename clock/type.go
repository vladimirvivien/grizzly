package clock

import (
	"time"
)

type Ticks = <-chan time.Time
type Clock = *clock

type clock struct {
	ticker *time.Ticker
}

func New(period time.Duration) Clock {
	return &clock{ticker:time.NewTicker(period)}
}

func (c *clock) Ticks() Ticks {
	return c.ticker.C
}

func (c *clock) Stop() {
	c.ticker.Stop()
}