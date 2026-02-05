// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build arm

package runtime

func checkgoarm() {
	// osinit not called yet, so numCPUStartup not set: must use
	// getCPUCount directly.
	if getCPUCount() > 1 && goarm < 7 {
		print("runtime: this system has multiple CPUs and must use\n")
		print("atomic synchronization instructions. Recompile using GOARM=7.\n")
		exit(1)
	}
}

//go:nosplit
func cputicks() int64 {
	// nanotime() is a poor approximation of CPU ticks that is enough for the profiler.
	return nanotime()
}
