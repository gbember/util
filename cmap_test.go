// cmap_test.go
package util

import "testing"

var (
	cm = NewCMap(900000)
)

func init() {
	for i := 0; i < 500000; i++ {
		cm.Put(i, i)
	}
}

func BenchmarkCMapPut(b *testing.B) {
	cm.Put(250000, 250000)
	cm.Put(250000, 250000)
	cm.Put(250000, 250000)
	cm.Put(250000, 250000)
	cm.Put(250000, 250000)
}
func BenchmarkCMapGet(b *testing.B) {
	cm.Get(250000)
	cm.Get(250000)
	cm.Get(250000)
	cm.Get(250000)
	cm.Get(250000)
}
