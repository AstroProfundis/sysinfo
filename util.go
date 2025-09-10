// Copyright © 2016 Zlatko Čalušić
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package sysinfo

import (
	"os"
	"strconv"
	"strings"
)

// Read one-liner text files, strip newline.
func slurpFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	// Trim spaces & \u0000 \uffff
	return strings.Trim(string(data), " \r\n\t\u0000\uffff")
}

// Write one-liner text files, add newline, ignore errors (best effort).
func spewFile(path string, data string, perm os.FileMode) {
	_ = os.WriteFile(path, []byte(data+"\n"), perm)
}

func SlurpFile(path string) string {
	return slurpFile(path)
}

func parseMemSize(key, memInfo string) uint64 {
	for _, line := range strings.Split(memInfo, "\n") {
		if !strings.Contains(line, key) {
			continue
		}
		fields := strings.Fields(line)
		size, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return 0
		}
		return size
	}
	return 0
}
