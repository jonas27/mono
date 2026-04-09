package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type procStats struct {
	MemRSS  uint64  // resident set size in bytes
	MemHeap uint64  // Go heap alloc in bytes
	CPU     float64 // CPU percentage since last sample

	prevTicks uint64
	prevTime  time.Time
}

func (s *procStats) update() {
	// Go heap from runtime
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	s.MemHeap = ms.HeapAlloc

	// RSS from /proc/self/status
	if data, err := os.ReadFile("/proc/self/status"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "VmRSS:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if kb, err := strconv.ParseUint(fields[1], 10, 64); err == nil {
						s.MemRSS = kb * 1024
					}
				}
				break
			}
		}
	}

	// CPU from /proc/self/stat (fields 14+15 are utime+stime in CLK_TCK ticks)
	if data, err := os.ReadFile("/proc/self/stat"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 15 {
			utime, _ := strconv.ParseUint(fields[13], 10, 64)
			stime, _ := strconv.ParseUint(fields[14], 10, 64)
			ticks := utime + stime
			now := time.Now()
			if !s.prevTime.IsZero() && ticks >= s.prevTicks {
				elapsed := now.Sub(s.prevTime).Seconds()
				if elapsed > 0 {
					// CLK_TCK = 100 on virtually all Linux systems
					s.CPU = float64(ticks-s.prevTicks) / (elapsed * 100) * 100
				}
			}
			s.prevTicks = ticks
			s.prevTime = now
		}
	}
}

func formatBytes(b uint64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/float64(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/float64(1<<20))
	default:
		return fmt.Sprintf("%.1f KB", float64(b)/float64(1<<10))
	}
}
