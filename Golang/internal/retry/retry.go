package retry

import (
	"time"
)

func Retry(attempts int, fn func() error) error {
	var err error
	delay := time.Second

	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}

		time.Sleep(delay)
		delay *= 2
	}

	return err
}
