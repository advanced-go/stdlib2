package controller

import "time"

type Timeout struct {
	DurationS string `json:"duration"`
	Duration  time.Duration
}

func NewTimeout(d time.Duration) *Timeout {
	t := new(Timeout)
	t.Duration = d
	return t
}
