package runtimecfg

import (
	"log"
	"runtime"
	"runtime/debug"
	"syscall"
	"unsafe"
)

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// Apply 根据当前机器的 CPU 核心数与内存情况调整 Go 运行时参数。
// 目标：让任务跑满可用 CPU，同时给界面留出调度空间，并避免多账号运行时内存无限增长。
func Apply() {
	cpu := runtime.NumCPU()
	if cpu < 1 {
		cpu = 1
	}
	runtime.GOMAXPROCS(cpu)

	gcPercent := 100
	if cpu >= 8 {
		gcPercent = 150
	}
	debug.SetGCPercent(gcPercent)

	total, avail := systemMemory()
	if total > 0 {
		limit := int64(float64(total) * 0.70)
		if avail > 0 {
			availLimit := int64(float64(avail) * 0.85)
			if availLimit < limit {
				limit = availLimit
			}
		}
		if limit < 512<<20 {
			limit = 512 << 20
		}
		debug.SetMemoryLimit(limit)
		log.Printf("运行环境: CPU=%d核(%s), 总内存=%.1fGB, 可用内存=%.1fGB, GC=%d%%, 内存软上限=%.1fGB",
			cpu, runtime.GOARCH, float64(total)/(1<<30), float64(avail)/(1<<30), gcPercent, float64(limit)/(1<<30))
		return
	}
	log.Printf("运行环境: CPU=%d核(%s), GC=%d%%", cpu, runtime.GOARCH, gcPercent)
}

func systemMemory() (total, avail uint64) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	proc := kernel32.NewProc("GlobalMemoryStatusEx")
	if proc.Find() != nil {
		return 0, 0
	}
	var status memoryStatusEx
	status.Length = uint32(unsafe.Sizeof(status))
	ret, _, _ := proc.Call(uintptr(unsafe.Pointer(&status)))
	if ret == 0 {
		return 0, 0
	}
	return status.TotalPhys, status.AvailPhys
}
