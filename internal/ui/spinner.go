package ui

import (
	"fmt"
	"os"
	"time"
)

// StopFunc stops a running spinner when called.
type StopFunc func()

// Spinner starts a non-blocking spinner writing to stderr and returns a stop function.
func Spinner(label string) StopFunc {
	done := make(chan struct{})
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		i := 0
		for {
			select {
			case <-done:
				fmt.Fprintf(os.Stderr, "\r%s\r", spaces(len(label)+4))
				return
			default:
				fmt.Fprintf(os.Stderr, "\r%s %s", frames[i%len(frames)], label)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return func() { close(done) }
}

func spaces(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}
