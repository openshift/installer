package smoke

import "fmt"

type NonCommented struct {
	numeric      int
	alphanumeric string
}

// timeout is a useful helper that returns a chan that is closed after
// the specified duration. This allows selecting the timeout to return
// immediately.
func timeout(d time.Duration) <-chan struct{} {
	t := time.NewTimer(d)
	c := make(chan struct{})
	go func() {
		<-t.C
		close(c)
	}()
	return c
}
