package retry

import "time"

type Strategy interface {
	Next() (time.Duration, bool)
}
