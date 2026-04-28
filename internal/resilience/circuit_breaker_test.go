package resilience

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)

	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)

	err := cb.Execute(func() error {
		return errors.New("some error")
	})

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestCircuitBreaker_Execute_ThresholdReached(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)

	for i := 0; i < 3; i++ {
		_ = cb.Execute(func() error {
			return errors.New("failure")
		})
	}

	state := cb.state
	if state != Open {
		t.Errorf("expected state Open, got %d", state)
	}
}

func TestCircuitBreaker_Execute_OpenStateRejected(t *testing.T) {
	cb := NewCircuitBreaker(2, time.Second)

	_ = cb.Execute(func() error {
		return errors.New("failure")
	})
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	err := cb.Execute(func() error {
		return nil
	})

	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_Execute_CooldownExpired(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	_ = cb.Execute(func() error {
		return errors.New("failure")
	})
	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	time.Sleep(100 * time.Millisecond)

	err := cb.Execute(func() error {
		return nil
	})

	if err != nil {
		t.Errorf("expected nil after cooldown, got %v", err)
	}

	state := cb.state
	if state != Closed {
		t.Errorf("expected state Closed, got %d", state)
	}
}

func TestCircuitBreaker_Execute_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(3, time.Second)

	_ = cb.Execute(func() error {
		return errors.New("failure")
	})

	_ = cb.Execute(func() error {
		return nil
	})

	failureCount := cb.failureCount
	if failureCount != 0 {
		t.Errorf("expected failureCount 0, got %d", failureCount)
	}
}