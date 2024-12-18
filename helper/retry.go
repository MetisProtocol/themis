package helper

import (
	"fmt"
	"time"
)

type stop struct {
	error
}

func NoRetryError(err error) stop {
	return stop{err}
}

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if s, ok := err.(stop); ok {
			return s.error
		}

		if attempts--; attempts > 0 {
			fmt.Printf("retry func error: %s. attempts #%d after %s.\n", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
		return err
	}
	return nil
}
