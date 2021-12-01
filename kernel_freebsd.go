// Copyright © 2016 Zlatko Čalušić
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.
//
//go:build freebsd
// +build freebsd

package sysinfo

import (
	"encoding/binary"
	"strconv"
	"syscall"
)

// Kernel information.
type Kernel struct {
	Release      string `json:"release,omitempty"`
	Version      string `json:"version,omitempty"`
	Architecture string `json:"architecture,omitempty"`
}

func (si *SysInfo) getKernelInfo() {
	osrel, err := syscall.Sysctl("kern.osrelease")
	if err != nil {
		return
	}
	si.Kernel.Release = osrel

	osver, err := syscall.Sysctl("kern.osrevision")
	if err != nil {
		return
	}
	osverInt := uint64(binary.LittleEndian.Uint32(append([]byte(osver), 0x0)))
	si.Kernel.Version = strconv.FormatUint(osverInt, 10)

	osarch, err := syscall.Sysctl("hw.machine_arch")
	if err != nil {
		return
	}
	si.Kernel.Architecture = osarch
}
