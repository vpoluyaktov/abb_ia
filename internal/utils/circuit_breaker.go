package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"abb_ia/internal/logger"
)

// CircuitBreakerState represents the current state of the circuit breaker
type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota // Normal operation
	StateOpen                              // Failing, not allowing requests
	StateHalfOpen                          // Testing if service is back
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	mu sync.RWMutex

	name string
	state CircuitBreakerState

	failureThreshold  int
	failureCount     int
	resetTimeout     time.Duration
	halfOpenTimeout  time.Duration
	lastStateChange  time.Time

	onStateChange func(name string, from, to CircuitBreakerState)
}

// CircuitBreakerOption defines options for creating a circuit breaker
type CircuitBreakerOption func(*CircuitBreaker)

// NewCircuitBreaker creates a new circuit breaker with the given options
func NewCircuitBreaker(name string, options ...CircuitBreakerOption) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:             name,
		state:            StateClosed,
		failureThreshold: 5,                // Default: 5 failures
		resetTimeout:     10 * time.Second, // Default: 10 seconds
		halfOpenTimeout:  2 * time.Second,  // Default: 2 seconds
		lastStateChange:  time.Now(),
		onStateChange: func(name string, from, to CircuitBreakerState) {
			logger.Info(fmt.Sprintf("Circuit breaker state changed - Name: %s, From: %s, To: %s",
				name, stateToString(from), stateToString(to)))
		},
	}

	// Apply options
	for _, option := range options {
		option(cb)
	}

	return cb
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.allowRequest() {
		return errors.New("circuit breaker is open")
	}

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}

	return err
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastStateChange) > cb.resetTimeout {
			cb.mu.RUnlock()
			cb.mu.Lock()
			cb.toHalfOpen()
			cb.mu.Unlock()
			cb.mu.RLock()
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

func (cb *CircuitBreaker) recordFailure() {
	cb.failureCount++
	
	if cb.state == StateHalfOpen || cb.failureCount >= cb.failureThreshold {
		cb.toOpen()
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	if cb.state == StateHalfOpen {
		cb.toClosed()
	}
	cb.failureCount = 0
}

func (cb *CircuitBreaker) toOpen() {
	if cb.state != StateOpen {
		oldState := cb.state
		cb.state = StateOpen
		cb.lastStateChange = time.Now()
		cb.onStateChange(cb.name, oldState, StateOpen)
	}
}

func (cb *CircuitBreaker) toHalfOpen() {
	if cb.state != StateHalfOpen {
		oldState := cb.state
		cb.state = StateHalfOpen
		cb.lastStateChange = time.Now()
		cb.onStateChange(cb.name, oldState, StateHalfOpen)
	}
}

func (cb *CircuitBreaker) toClosed() {
	if cb.state != StateClosed {
		oldState := cb.state
		cb.state = StateClosed
		cb.lastStateChange = time.Now()
		cb.failureCount = 0
		cb.onStateChange(cb.name, oldState, StateClosed)
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitBreakerState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// WithFailureThreshold sets the failure threshold
func WithFailureThreshold(threshold int) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.failureThreshold = threshold
	}
}

// WithResetTimeout sets the reset timeout
func WithResetTimeout(timeout time.Duration) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.resetTimeout = timeout
	}
}

// WithHalfOpenTimeout sets the half-open timeout
func WithHalfOpenTimeout(timeout time.Duration) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.halfOpenTimeout = timeout
	}
}

// WithStateChangeHandler sets the state change handler
func WithStateChangeHandler(handler func(name string, from, to CircuitBreakerState)) CircuitBreakerOption {
	return func(cb *CircuitBreaker) {
		cb.onStateChange = handler
	}
}

func stateToString(state CircuitBreakerState) string {
	switch state {
	case StateClosed:
		return "CLOSED"
	case StateOpen:
		return "OPEN"
	case StateHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}
