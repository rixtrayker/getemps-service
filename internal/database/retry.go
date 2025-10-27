package database

import (
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/sirupsen/logrus"
)

func ExecuteWithRetry(operation func() error, logger *logrus.Logger) error {
	return retry.Do(
		operation,
		retry.Attempts(3),
		retry.Delay(100*time.Millisecond),
		retry.DelayType(retry.BackOffDelay),
		retry.OnRetry(func(n uint, err error) {
			if logger != nil {
				logger.WithFields(logrus.Fields{
					"attempt": n + 1,
					"error":   err.Error(),
				}).Warn("Database operation retry")
			}
		}),
		retry.RetryIf(func(err error) bool {
			// Retry on connection errors, timeout errors, etc.
			// Add specific error checking here based on your needs
			return err != nil
		}),
	)
}
