// Copyright © 2016 Zlatko Čalušić
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package sysinfo

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// NetworkDevice information.
type NetworkDevice struct {
	Name       string   `json:"name,omitempty"`
	Driver     string   `json:"driver,omitempty"`
	MACAddress string   `json:"macaddress,omitempty"`
	IPAddress  []string `json:"ipaddress,omitempty"`
	Port       string   `json:"port,omitempty"`
	Speed      uint     `json:"speed,omitempty"` // device max supported speed in Mbps
	MTU        int      `json:"mtu,omitempty"`
}

func getPortType(supp uint32) (port string) {
	for i, p := range [...]string{"tp", "aui", "mii", "fibre", "bnc"} {
		if supp&(1<<uint(i+7)) > 0 {
			port += p + "/"
		}
	}

	port = strings.TrimRight(port, "/")
	return
}

func getMaxSpeed(supp uint32) (speed uint) {
	// Fancy, right?
	switch {
	case supp&0x78000000 > 0:
		speed = 56000
	case supp&0x07800000 > 0:
		speed = 40000
	case supp&0x00600000 > 0:
		speed = 20000
	case supp&0x001c1000 > 0:
		speed = 10000
	case supp&0x00008000 > 0:
		speed = 2500
	case supp&0x00020030 > 0:
		speed = 1000
	case supp&0x0000000c > 0:
		speed = 100
	case supp&0x00000003 > 0:
		speed = 10
	}

	return
}

func getSupported(name string) uint32 {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return 0
	}
	defer syscall.Close(fd)

	// struct ethtool_cmd from /usr/include/linux/ethtool.h
	var ethtool struct {
		Cmd       uint32
		Supported uint32
	}

	// ETHTOOL_GSET from /usr/include/linux/ethtool.h
	const GSET = 0x1

	ethtool.Cmd = GSET

	// struct ifreq from /usr/include/linux/if.h
	var ifr struct {
		Name [16]byte
		Data uintptr
	}

	copy(ifr.Name[:], name+"\000")
	ifr.Data = uintptr(unsafe.Pointer(&ethtool))

	// SIOCETHTOOL from /usr/include/linux/sockios.h
	const SIOCETHTOOL = 0x8946

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(SIOCETHTOOL), uintptr(unsafe.Pointer(&ifr)))
	if errno == 0 {
		return ethtool.Supported
	}

	return 0
}

// readIfaceSpeed reads /sys/class/net/<iface>/speed
// inspired by https://github.com/prometheus/procfs/blob/master/sysfs/net_class.go#L56
func readIfaceSpeed(file string) (speed uint) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	if b, err := ioutil.ReadAll(f); err == nil {
		if s, err := strconv.Atoi(strings.TrimSpace(string(b))); err == nil && s > 0 {
			speed = uint(s)
		}
	}
	return
}

func (si *SysInfo) getNetworkInfo() {
	sysClassNet := "/sys/class/net"
	devices, err := ioutil.ReadDir(sysClassNet)
	if err != nil {
		return
	}

	si.Network = make([]NetworkDevice, 0)
	for _, link := range devices {
		fullpath := path.Join(sysClassNet, link.Name())
		dev, err := os.Readlink(fullpath)
		if err != nil {
			continue
		}

		if strings.HasPrefix(dev, "../../devices/virtual/net/lo") ||
			strings.HasPrefix(dev, "../../devices/virtual/net/dummy") {
			continue
		}

		supp := getSupported(link.Name())

		// get IP Address(es) of interface
		inet, _ := net.InterfaceByName(link.Name())
		ip_addrs, _ := getIPAddrByInterface(inet)

		speed := readIfaceSpeed(path.Join(fullpath, "speed"))
		if speed == 0 {
			speed = getMaxSpeed(supp)
		}

		device := NetworkDevice{
			Name:       link.Name(),
			MACAddress: slurpFile(path.Join(fullpath, "address")),
			IPAddress:  ip_addrs,
			Port:       getPortType(supp),
			Speed:      speed,
			MTU:        inet.MTU,
		}

		if driver, err := os.Readlink(path.Join(fullpath, "device", "driver")); err == nil {
			device.Driver = path.Base(driver)
		}

		si.Network = append(si.Network, device)
	}
}

func getIPAddrByInterface(inet *net.Interface) ([]string, error) {
	ip_addrs := make([]string, 0)
	addrs, err := inet.Addrs()
	if err != nil {
		return ip_addrs, err
	}
	for _, addr := range addrs {
		ip_addrs = append(ip_addrs, addr.String())
	}
	return ip_addrs, err
}
