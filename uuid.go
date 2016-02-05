// uuid.go
package util

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

var (
	guuid = NewUUIDCreater()
)

//UUID生成器
type UUIDCreater struct {
	machineIds []byte
	counter    uint32
}

func UUID() uint64 {
	return guuid.Id()
}

//生成ID生成器
func NewUUIDCreater() *UUIDCreater {
	uuid := &UUIDCreater{}
	uuid.counter = readRandomUint32()
	uuid.machineIds = readMachineId()
	return uuid
}

func (uuid *UUIDCreater) Id() uint64 {
	var b [12]byte
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	b[4] = uuid.machineIds[0]
	b[5] = uuid.machineIds[1]
	b[6] = uuid.machineIds[2]
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	i := atomic.AddUint32(&uuid.counter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	b1 := b[:]
	return (binary.BigEndian.Uint64(b1[:8]) << 3) + uint64(binary.BigEndian.Uint32(b1[8:]))
}
func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot read random object id: %v", err))
	}
	return uint32((uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}
func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(rand.Reader, id)
		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}
