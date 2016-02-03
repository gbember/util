// gtimer_test.go
package util

import (
	"testing"
	"time"
)

func TestGTimer(t *testing.T) {
	gtmr := NewGTimer()
	//	time.Sleep(time.Second)
	gtmr.Stop()
	time.Sleep(time.Second * 5)
}
