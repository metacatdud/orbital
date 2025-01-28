package machine

import (
	"errors"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

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
	Total string `json:"total"`
	Free  string `json:"free"`
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

func getCPUInfo() (*CPUInfo, error) {
	cpuInfo := &CPUInfo{}

	if cores, err := cpu.Counts(true); err == nil {
		cpuInfo.TotalCores = cores
	}

	if loads, err := cpu.Percent(1*time.Second, true); err == nil {
		cpuInfo.CoreLoad = loads
	}

	return cpuInfo, nil
}

func getMemInfo() (*MemoryInfo, error) {
	memInfo := &MemoryInfo{}

	if vm, err := mem.VirtualMemory(); err == nil {
		memInfo.Total = vm.Total
		memInfo.Free = vm.Available
	}

	return memInfo, nil
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
