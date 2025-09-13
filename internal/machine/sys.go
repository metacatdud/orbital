package machine

import (
	"errors"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

type Info struct {
	HostID          string
	Hostname        string
	Uptime          uint64
	Procs           uint64
	OS              string
	Platform        string
	PlatformFamily  string
	PlatformVersion string
	KernelVersion   string
	KernelArch      string
}

type CPUInfo struct {
	TotalCores int       `json:"totalCores"`
	CoreLoad   []float64 `json:"coreLoad"`
}

type GPU struct {
	Model  string `json:"model"`
	Vendor string `json:"vendor"`
}

type GPUInfo struct {
	GPUs []GPU `json:"gpus"`
}

type Disk struct {
	Name  string `json:"name"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

type DiskInfo struct {
	Total uint64 `json:"total"`
	Disks []Disk `json:"disks"`
}

type MemoryInfo struct {
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

type NetworkInfo struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func getInfo() (*Info, error) {
	sysInfo := &Info{}

	info, err := host.Info()
	if err != nil {
		return nil, err
	}

	sysInfo.HostID = info.HostID
	sysInfo.Hostname = info.Hostname
	sysInfo.Uptime = info.Uptime
	sysInfo.Procs = info.Procs
	sysInfo.Platform = info.Platform
	sysInfo.PlatformFamily = info.PlatformFamily
	sysInfo.PlatformVersion = info.PlatformVersion
	sysInfo.KernelVersion = info.KernelVersion
	sysInfo.KernelArch = info.KernelArch

	return sysInfo, nil

}

func getCPUInfo() (*CPUInfo, error) {
	info := &CPUInfo{}

	if cores, err := cpu.Counts(true); err == nil {
		info.TotalCores = cores
	}

	if loads, err := cpu.Percent(1*time.Second, true); err == nil {
		info.CoreLoad = loads
	}

	return info, nil
}

func getDiskInfo() (*DiskInfo, error) {
	info := &DiskInfo{}

	partitions, err := disk.Partitions(false)
	if err != nil {
		return info, err
	}

	for _, p := range partitions {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			// TODO: Maybe log this errors somehow?
			continue // skip partitions with issues
		}

		d := Disk{
			Name:  p.Device,
			Total: usage.Total,
			Free:  usage.Free,
		}

		info.Disks = append(info.Disks, d)
		info.Total += usage.Total
	}

	return info, nil
}

func getMemInfo() (*MemoryInfo, error) {
	info := &MemoryInfo{}

	if vm, err := mem.VirtualMemory(); err == nil {
		info.Total = vm.Total
		info.Free = vm.Available
	}

	return info, nil
}

func getNetworkInfo() (*NetworkInfo, error) {

	iface, err := getDefaultNetworkInterface()
	if err != nil {
		return &NetworkInfo{
			Type: "unknown",
			Name: "",
		}, nil
	}

	netType := "cable"
	if isWirelessInterface(iface) {
		netType = "wifi"
	}

	name := "wired"
	if netType == "wifi" {
		name = getSSID(iface)
	}

	return &NetworkInfo{
		Type: netType,
		Name: name,
	}, nil
}

func getDefaultNetworkInterface() (string, error) {
	stats, err := net.IOCounters(true)
	if err != nil {
		return "", err
	}

	var (
		maxBytes     uint64
		defaultIface string
	)

	for _, s := range stats {
		if s.BytesSent > maxBytes && s.Name != "lo" {
			maxBytes = s.BytesSent
			defaultIface = s.Name
		}
	}

	if defaultIface == "" {
		return "", errors.New("no default network interface")
	}

	return defaultIface, nil
}

func getSSID(iface string) string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxSSID(iface)
	}
	return "unknown"
}

func isWirelessInterface(iface string) bool {
	switch runtime.GOOS {
	case "linux":
		return strings.HasPrefix(iface, "wl") || strings.Contains(iface, "wlan")
	}

	return false
}

func getLinuxSSID(iface string) string {
	cmd := exec.Command("iwgetid", "-r")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}
