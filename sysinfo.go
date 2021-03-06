// Copyright © 2016 Zlatko Čalušić
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

// Package sysinfo is a Go library providing Linux OS / kernel / hardware system information.
package sysinfo

// SysInfo struct encapsulates all other information structs.
type SysInfo struct {
	Node    Node            `json:"node,omitempty"`
	OS      OS              `json:"os,omitempty"`
	Kernel  Kernel          `json:"kernel,omitempty"`
	Product Product         `json:"product,omitempty"`
	Board   Board           `json:"board,omitempty"`
	Chassis Chassis         `json:"chassis,omitempty"`
	BIOS    BIOS            `json:"bios,omitempty"`
	CPU     CPU             `json:"cpu,omitempty"`
	Memory  Memory          `json:"memory,omitempty"`
	Storage []StorageDevice `json:"storage,omitempty"`
	LVM     []LogicalVolume `json:"lvm,omitempty"`
	Network []NetworkDevice `json:"network,omitempty"`
}

// GetSysInfo gathers all available system information.
func (si *SysInfo) GetSysInfo() {
	// DMI info
	si.getProductInfo()
	si.getBoardInfo()
	si.getChassisInfo()
	si.getBIOSInfo()

	// SMBIOS info
	si.getMemoryInfo()

	// Node info
	si.getNodeInfo() // depends on BIOS info

	// Hardware info
	si.getCPUInfo() // depends on Node info
	si.getStorageInfo()
	si.getLVMInfo()
	si.getNetworkInfo()

	// Software info
	si.getOSInfo()
	si.getKernelInfo()
}
