package domain

import (
	"encoding/json"
	"fmt"
)

type Disk struct {
	Name  string `json:"name"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

type DiskInfo struct {
	Total uint64 `json:"total"`
	Disks []Disk `json:"disks"`
}

type Machine struct {
	HostID          string   `json:"hostId"`
	Hostname        string   `json:"hostname"`
	Uptime          uint64   `json:"uptime"`
	Procs           uint64   `json:"procs"`
	OS              string   `json:"os"`
	Platform        string   `json:"platform"`
	PlatformFamily  string   `json:"platformFamily"`
	PlatformVersion string   `json:"platformVersion"`
	KernelVersion   string   `json:"kernelVersion"`
	KernelArch      string   `json:"kernelArch"`
	TotalCores      int      `json:"totalCores"`
	CoreLoads       []string `json:"coreLoads"`   // Each load as "XX.XX%"
	TotalMemGB      string   `json:"totalMemGB"`  // Total memory in GB, formatted to 2 decimals
	FreeMemGB       string   `json:"freeMemGB"`   // Free memory in GB, formatted to 2 decimals
	FreePercent     string   `json:"freePercent"` // Percentage of free memory, formatted to 2 decimals
	UsedPercent     string   `json:"usedPercent"`
	NetType         string   `json:"netType"`
	NetName         string   `json:"netName"`
	Disks           DiskInfo `json:"disks"`
}

type Machines []Machine

func NewMachineFromData(machineBin []byte) (*Machine, error) {
	var wrapper struct {
		SystemInfo rawMachine `json:"systemInfo"`
	}

	if err := json.Unmarshal(machineBin, &wrapper); err != nil {
		return nil, err
	}

	return createMachineInfoFromRaw(wrapper.SystemInfo), nil
}

func createMachineInfoFromRaw(raw rawMachine) *Machine {
	var m = &Machine{}

	m.HostID = raw.Info.HostID
	m.Hostname = raw.Info.Hostname
	m.Uptime = raw.Info.Uptime
	m.Procs = raw.Info.Procs
	m.OS = raw.Info.OS
	m.Platform = raw.Info.Platform
	m.PlatformFamily = raw.Info.PlatformFamily
	m.PlatformVersion = raw.Info.PlatformVersion
	m.KernelVersion = raw.Info.KernelVersion
	m.KernelArch = raw.Info.KernelArch

	m.TotalCores = raw.CPU.TotalCores

	m.CoreLoads = make([]string, len(raw.CPU.CoreLoad))
	for i, load := range raw.CPU.CoreLoad {
		m.CoreLoads[i] = fmt.Sprintf("%.2f%%", load)
	}

	totalGB := float64(raw.Mem.Total) / (1024 * 1024 * 1024)
	freeGB := float64(raw.Mem.Free) / (1024 * 1024 * 1024)
	m.TotalMemGB = fmt.Sprintf("%.2f", totalGB)
	m.FreeMemGB = fmt.Sprintf("%.2f", freeGB)

	if raw.Mem.Total > 0 {
		freeRatio := (float64(raw.Mem.Free) / float64(raw.Mem.Total)) * 100
		usedRatio := 100 - freeRatio
		m.FreePercent = fmt.Sprintf("%.2f", freeRatio)
		m.UsedPercent = fmt.Sprintf("%.2f", usedRatio)
	}

	m.NetType = raw.Net.Type
	m.NetName = raw.Net.Name

	m.Disks = raw.Disks

	return m
}

type rawMachine struct {
	Info struct {
		HostID          string `json:"hostId"`
		Hostname        string `json:"hostname"`
		Uptime          uint64 `json:"uptime"`
		Procs           uint64 `json:"procs"`
		OS              string `json:"os"`
		Platform        string `json:"platform"`
		PlatformFamily  string `json:"platformFamily"`
		PlatformVersion string `json:"platformVersion"`
		KernelVersion   string `json:"kernelVersion"`
		KernelArch      string `json:"kernelArch"`
	} `json:"info"`

	Disks DiskInfo `json:"disks"`

	CPU struct {
		TotalCores int       `json:"totalCores"`
		CoreLoad   []float64 `json:"coreLoad"`
	} `json:"cpu"`

	Mem struct {
		Total uint64 `json:"total"`
		Free  uint64 `json:"free"`
	} `json:"mem"`

	Net struct {
		Type string `json:"type"`
		Name string `json:"name"`
	} `json:"net"`
}
