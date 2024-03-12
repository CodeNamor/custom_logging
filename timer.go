package logging

import "time"

type Timer struct {
	start, end time.Time
}

func (t *Timer) Start() {
	t.start = time.Now()
	t.end = time.Time{}
}

func (t *Timer) Stop() {
	t.end = time.Now()
}

func (t *Timer) Elapsed() time.Duration {
	if t.start.IsZero() {
		return 0
	}

	if t.end.IsZero() {
		return time.Since(t.start) / time.Millisecond
	}

	return t.end.Sub(t.start) / time.Millisecond
}

func (t *Timer) GetStartTime() time.Time {
	return t.start
}

func (t *Timer) GetEndTime() time.Time {
	return t.end
}
