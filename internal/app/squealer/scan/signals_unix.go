// +build linux bsd darwin

package scan

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func (s *gitScanner) monitorSignals(processes int, wg sync.WaitGroup) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGTSTP)
	go func() {
		for _ = range c {
			fmt.Println("\r- Exiting")
			for i := 0; i < processes; i++ {
				wg.Done()
			}
			os.Exit(0)
		}
	}()
}
