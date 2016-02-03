package util

import (
	"math"
	"runtime"
	"time"
)

type GCStats struct {
	Time          int64   `json:"t"`      // 时间
	HeapAlloc     uint64  `json:"ha"`     // 已经分配的内存
	HeapObjects   uint64  `json:"hos"`    // 分配的对象总数量
	GCMaxPause    uint64  `json:"gcmax"`  // gc最大暂停时间
	GCAvgPause    uint64  `json:"gcavg"`  // gc平均暂停时间
	GCRuns        uint32  `json:"gcn"`    // 距离上一次统计gc执行次数
	GCCPUFraction float64 `json:"gccpuf"` // gc消耗cpu时间
}

var (
	lastMemStats = &runtime.MemStats{}
)

//得到GC信息
func GetGCState() *GCStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// sort the GC pause array
	length := len(memStats.PauseNs)
	if int(memStats.NumGC) < length {
		length = int(memStats.NumGC)
	}

	maxPause := uint64(0)
	pause := uint64(0)
	for i := lastMemStats.NumGC + 1; i <= memStats.NumGC; i++ {
		pause = memStats.PauseNs[(i+255)%256]
		if pause > maxPause {
			maxPause = pause
		}

	}
	gcNum := memStats.NumGC - lastMemStats.NumGC

	ret := &GCStats{
		Time:          time.Now().UnixNano(),
		HeapAlloc:     memStats.HeapAlloc,
		HeapObjects:   memStats.HeapObjects,
		GCMaxPause:    maxPause,
		GCAvgPause:    (memStats.PauseTotalNs - lastMemStats.PauseTotalNs) / uint64(gcNum),
		GCRuns:        gcNum,
		GCCPUFraction: memStats.GCCPUFraction,
	}

	lastMemStats = &memStats

	return ret
}

func percentile(perc float64, arr []uint64, length int) uint64 {
	if length == 0 {
		return 0
	}
	indexOfPerc := int(math.Floor(((perc / 100.0) * float64(length)) + 0.5))
	if indexOfPerc >= length {
		indexOfPerc = length - 1
	}
	return arr[indexOfPerc]
}
