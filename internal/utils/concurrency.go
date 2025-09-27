package utils

import (
	"context"
	"os"
	"strconv"
	"sync"
)

const (
	defaultMaxConcurrentRequests = 8
	maxConcurrentEnvVar          = "RIPESTAT_MAX_CONCURRENCY"
)

var (
	maxConcurrentRequests = defaultMaxConcurrentRequests
	maxConcurrentMu       sync.RWMutex
)

func init() {
	if val := os.Getenv(maxConcurrentEnvVar); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			if parsed > defaultMaxConcurrentRequests {
				parsed = defaultMaxConcurrentRequests
			}
			maxConcurrentRequests = parsed
		}
	}
}

func SetMaxConcurrentRequests(n int) {
	if n <= 0 {
		n = 1
	}
	if n > defaultMaxConcurrentRequests {
		n = defaultMaxConcurrentRequests
	}

	maxConcurrentMu.Lock()
	maxConcurrentRequests = n
	maxConcurrentMu.Unlock()
}

func GetMaxConcurrentRequests() int {
	maxConcurrentMu.RLock()
	defer maxConcurrentMu.RUnlock()
	return maxConcurrentRequests
}

func runWithConcurrency[T any, R any](ctx context.Context, items []T, fn func(context.Context, int, T) (R, error)) ([]R, error) {
	maxConcurrentMu.RLock()
	maxWorkers := maxConcurrentRequests
	maxConcurrentMu.RUnlock()

	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]R, len(items))
	sem := make(chan struct{}, maxWorkers)
	errCh := make(chan error, len(items))

	var wg sync.WaitGroup

	for idx, item := range items {
		idx := idx
		item := item

		wg.Add(1)
		go func() {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			result, err := fn(ctx, idx, item)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				cancel()
				return
			}

			results[idx] = result
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	if err := ctx.Err(); err != nil && len(items) > 0 {
		return nil, err
	}

	return results, nil
}
