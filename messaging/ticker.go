package messaging

import "time"

type Ticker struct {
	duration time.Duration
	original time.Duration
	ticker   *time.Ticker
}

func NewTicker(duration time.Duration) *Ticker {
	t := new(Ticker)
	t.duration = duration
	t.original = duration
	return t
}

func (t *Ticker) Duration() time.Duration { return t.duration }
func (t *Ticker) C() <-chan time.Time     { return t.ticker.C }

func (t *Ticker) Start(newDuration time.Duration) {
	if newDuration <= 0 {
		newDuration = t.duration
	} else {
		t.duration = newDuration
	}
	t.Stop()
	t.ticker.Reset(newDuration)
}

func (t *Ticker) Reset() {
	t.Start(t.original)
}

func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
}
