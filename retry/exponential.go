package retry

import (
	"github.com/RyoLena/Gadget/internal/errs"
	"math"
	"sync/atomic"
	"time"
)

type ExponentialBackoffRetryStrategy struct {
	//初始重试间隔
	initialInterval time.Duration
	//最大重试次数
	maxInterval time.Duration
	//最大重试间隔
	maxRetries int32
	// 当前重试次数
	retries int32
	//是否已经达到最大重试间隔
	maxIntervalReached atomic.Value
}

func NewExponentialBackoffRetryStrategy(initialInterval, maxInterval time.Duration,
	maxRetries int32) (*ExponentialBackoffRetryStrategy, error) {
	if initialInterval <= 0 {
		return nil, errs.NewErrInvalidIntervalValue(initialInterval)
	} else if maxInterval < initialInterval {
		return nil, errs.NewErrInvalidMaxIntervalValue(maxInterval, initialInterval)
	}
	return &ExponentialBackoffRetryStrategy{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxRetries:      maxRetries,
	}, nil
}

func (s *ExponentialBackoffRetryStrategy) Next() (time.Duration, bool) {
	retries := atomic.AddInt32(&s.retries, 1)
	if s.maxRetries <= 0 || retries <= s.maxRetries {
		if reached, ok := s.maxIntervalReached.Load().(bool); ok && reached {
			return s.maxInterval, true
		}
		interval := s.initialInterval * time.Duration(math.Pow(2, float64(retries-1)))
		// 溢出或当前重试间隔大于最大重试间隔
		if interval <= 0 || interval > s.maxInterval {
			s.maxIntervalReached.Store(true)
			return s.maxInterval, true
		}
		return interval, true
	}
	return 0, false
}
