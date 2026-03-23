package resilience

import (
	"errors"
	"sync"
	"time"
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type State int

const (
	Closed State = iota
	Open
)

type CircuitBreaker struct {
	mu           sync.RWMutex
	state        State
	failureCount int
	threshold    int
	cooldown     time.Duration
	lastFailure  time.Time
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		cooldown:  cooldown,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	cb.mu.Lock()
	if cb.state == Open && time.Since(cb.lastFailure) > cb.cooldown {
		cb.state = Closed
		cb.failureCount = 0
	}

	if cb.state == Open {
		cb.mu.Unlock()
		return ErrCircuitOpen
	}
	cb.mu.Unlock()

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.failureCount++
		cb.lastFailure = time.Now()
		if cb.failureCount >= cb.threshold {
			cb.state = Open
		}
		return err
	}

	cb.failureCount = 0
	return nil
}
