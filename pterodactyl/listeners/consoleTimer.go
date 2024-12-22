package listeners

import "time"

type consoleTimer struct {
	t       *time.Timer
	started bool
	maxTime time.Duration
	c       <-chan time.Time
}

func newTimer(maxTime string) (*consoleTimer, error) {
	if maxTime == "" {
		return &consoleTimer{}, nil
	}
	m, err := time.ParseDuration(maxTime)
	if err != nil {
		return nil, err
	}
	return &consoleTimer{
		maxTime: m,
	}, nil
}

func (t *consoleTimer) start() {
	if t.started {
		return
	}
	t.t = time.NewTimer(t.maxTime)
	t.c = t.t.C
	t.started = true
}
func (t *consoleTimer) reset() {
	t.t.Reset(t.maxTime)
}
func (t *consoleTimer) stop() {
	t.t.Stop()
	t.started = false
}
