package retry

import (
	"log"
	"time"
)

type Func func() error

type config struct {
	attempts      int
	sleep         time.Duration
	backoff       bool
	backoffFactor float64
}

type Option func(*config)

func WithAttempts(attempts int) Option {
	return func(c *config) {
		c.attempts = attempts
	}
}

func WithBackoff(duration time.Duration, factor float64) Option {
	return func(c *config) {
		c.backoff = true
		c.sleep = duration
		c.backoffFactor = factor
	}
}

func Do(fn Func, opts ...Option) error {
	conf := config{attempts: 1, sleep: 1, backoff: false, backoffFactor: 1.0}
	for _, opt := range opts {
		opt(&conf)
	}

	var err error
	for i := range conf.attempts {
		if i > 0 {
			log.Printf("Retrying... attempt %d/%d", i, conf.attempts)
			time.Sleep(conf.sleep)
			if conf.backoff {
				conf.sleep = time.Duration(float64(conf.sleep) * conf.backoffFactor)
			}
		}
		err = fn()
		if err == nil {
			return nil
		}
	}
	return err
}
