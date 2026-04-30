package shutdown

import (
	"context"
	"log"
	"sync"
	"time"
)

// Manager handles graceful shutdown coordination across the application.
type Manager struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	shutdownCh chan struct{}
	once       sync.Once
}

// NewManager creates a new shutdown manager.
func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		ctx:        ctx,
		cancel:     cancel,
		shutdownCh: make(chan struct{}),
	}
}

// Context returns the shutdown context.
// Background workers should use this context to know when to stop.
func (m *Manager) Context() context.Context {
	return m.ctx
}

// AddWorker increments the worker count.
// Call this before starting a background goroutine.
func (m *Manager) AddWorker() {
	m.wg.Add(1)
}

// WorkerDone decrements the worker count.
// Call this when a background goroutine completes (typically in defer).
func (m *Manager) WorkerDone() {
	m.wg.Done()
}

// Shutdown initiates graceful shutdown.
// It can be called multiple times safely (only first call takes effect).
func (m *Manager) Shutdown() {
	m.once.Do(func() {
		log.Println("Shutdown manager: Initiating graceful shutdown...")
		close(m.shutdownCh)
		m.cancel() // Signal all workers to stop
	})
}

// Wait waits for all workers to complete with a timeout.
// Returns true if all workers completed, false if timeout occurred.
func (m *Manager) Wait(timeout time.Duration) bool {
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("Shutdown manager: All workers completed successfully")
		return true
	case <-time.After(timeout):
		log.Printf("Shutdown manager: Timeout after %v, some workers may not have completed", timeout)
		return false
	}
}

// ShutdownChannel returns a channel that closes when shutdown is initiated.
// Useful for select statements in workers.
func (m *Manager) ShutdownChannel() <-chan struct{} {
	return m.shutdownCh
}

// IsShuttingDown returns true if shutdown has been initiated.
func (m *Manager) IsShuttingDown() bool {
	select {
	case <-m.shutdownCh:
		return true
	default:
		return false
	}
}
