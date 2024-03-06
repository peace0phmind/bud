package info

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

type CpuCountInfo struct {
	Physical int `json:"physical"`
	Logical  int `json:"logical"`
}

type SysInfo struct {
}

func CpuCounts() (ret CpuCountInfo, err error) {
	ret.Physical, err = cpu.Counts(false)
	if err != nil {
		return
	}
	ret.Logical, err = cpu.Counts(true)
	return
}

func CpuTimes() ([]cpu.TimesStat, error) {
	return cpu.Times(true)
}

func Cpu() ([]cpu.InfoStat, error) {
	return cpu.Info()
}

func CpuPercent() ([]float64, error) {
	return cpu.Percent(1*time.Second, false)
}

func Memory() (*mem.VirtualMemoryStat, error) {
	return mem.VirtualMemory()
}

func Host() (*host.InfoStat, error) {
	return host.Info()
}

func LoadAvg() (*load.AvgStat, error) {
	return load.Avg()
}

func LoadMisc() (*load.MiscStat, error) {
	return load.Misc()
}

func NetInterfaces() (net.InterfaceStatList, error) {
	return net.Interfaces()
}

func NetConnections() ([]net.ConnectionStat, error) {
	return net.Connections("all")
}

func NetIOCounters() ([]net.IOCountersStat, error) {
	return net.IOCounters(true)
}

func DiskPartitions() ([]disk.PartitionStat, error) {
	return disk.Partitions(false)
}

func DiskUsage() ([]disk.UsageStat, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskUsageStats []disk.UsageStat
	for _, partition := range partitions {
		usage, err1 := disk.Usage(partition.Mountpoint)
		if err1 != nil {
			return nil, err1
		}
		diskUsageStats = append(diskUsageStats, *usage)
	}
	return diskUsageStats, nil
}
