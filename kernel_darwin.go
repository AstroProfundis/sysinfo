// Copyright © 2016 Zlatko Čalušić
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

//go:build darwin
// +build darwin

package sysinfo

// Kernel information.
type Kernel struct {
	Release      string `json:"release,omitempty"`
	Version      string `json:"version,omitempty"`
	Architecture string `json:"architecture,omitempty"`
}

func (si *SysInfo) getKernelInfo() {
	si.Kernel.Release = slurpFile("/proc/sys/kernel/osrelease")
	si.Kernel.Version = slurpFile("/proc/sys/kernel/version")

	/*
		// Uname() and Utsname is not defined in syscall for darwin
		var uname syscall.Utsname
		if err := syscall.Uname(&uname); err != nil {
			return
		}

		si.Kernel.Architecture = strings.TrimRight(string((*[65]byte)(unsafe.Pointer(&uname.Machine))[:]), "\000")
	*/
}
