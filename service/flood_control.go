package service

import (
	"context"
	"sync"
	"time"
)

type FloodControl interface {
	Check(ctx context.Context, userID int64) (bool, error)
}

type MemFloodControl struct {
	mtx         sync.Mutex
	requests    map[int64][]time.Time
	timeRange   time.Duration
	maxRequests int
}

func NewMemFloodControl(timeRange time.Duration, maxRequests int) *MemFloodControl {
	return &MemFloodControl{
		requests:    make(map[int64][]time.Time),
		timeRange:   timeRange,
		maxRequests: maxRequests,
	}
}

func (mfc *MemFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	mfc.mtx.Lock()
	defer mfc.mtx.Unlock()

	timeNow := time.Now()
	requests, ok := mfc.requests[userID]
	if !ok {
		requests = []time.Time{}
	}

	timeInRange := []time.Time{}
	for _, t := range requests {
		if timeNow.Sub(t) <= mfc.timeRange {
			timeInRange = append(timeInRange, t)
		}
	}

	timeInRange = append(timeInRange, timeNow)
	mfc.requests[userID] = timeInRange

	if len(timeInRange) > mfc.maxRequests {
		return false, nil
	}

	return true, nil
}
