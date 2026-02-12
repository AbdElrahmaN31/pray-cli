// Package ui provides interactive user interface components
package ui

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Spinner provides a simple terminal spinner for long-running operations
type Spinner struct {
	message  string
	frames   []string
	interval time.Duration
	mu       sync.Mutex
	running  bool
	done     chan struct{}
}

// NewSpinner creates a new spinner with the given message
func NewSpinner(message string) *Spinner {
	return &Spinner{
		message:  message,
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 80 * time.Millisecond,
		done:     make(chan struct{}),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		frameIdx := 0
		for {
			select {
			case <-s.done:
				return
			case <-ticker.C:
				s.mu.Lock()
				if !s.running {
					s.mu.Unlock()
					return
				}
				// Clear line and print spinner
				fmt.Fprintf(os.Stdout, "\r\033[K%s %s", s.frames[frameIdx], s.message)
				frameIdx = (frameIdx + 1) % len(s.frames)
				s.mu.Unlock()
			}
		}
	}()
}

// Stop stops the spinner without any message
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	close(s.done)
	fmt.Fprint(os.Stdout, "\r\033[K") // Clear line
}

// Success stops the spinner and shows a success message
func (s *Spinner) Success(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	close(s.done)

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Fprintf(os.Stdout, "\r\033[K%s %s\n", green("✓"), message)
}

// Fail stops the spinner and shows a failure message
func (s *Spinner) Fail(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	s.running = false
	close(s.done)

	red := color.New(color.FgRed).SprintFunc()
	fmt.Fprintf(os.Stdout, "\r\033[K%s %s\n", red("✗"), message)
}

// Update updates the spinner message while running
func (s *Spinner) Update(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}
