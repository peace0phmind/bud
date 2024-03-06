package info

import (
	"encoding/json"
	"fmt"
	"testing"
)

func printInfo(info any, err error) {
	if err != nil {
		return
	}
	b, err := json.MarshalIndent(info, "", "  ")
	fmt.Print(string(b))
}

func TestCpuCounts(t *testing.T) {
	printInfo(CpuCounts())
}

func TestCpuTimes(t *testing.T) {
	printInfo(CpuTimes())
}

func TestCpu(t *testing.T) {
	printInfo(Cpu())
}

func TestCpuPercent(t *testing.T) {
	printInfo(CpuPercent())
}

func TestMemory(t *testing.T) {
	printInfo(Memory())
}

func TestHost(t *testing.T) {
	printInfo(Host())
}

func TestLoadAvg(t *testing.T) {
	printInfo(LoadAvg())
}

func TestLoadMisc(t *testing.T) {
	printInfo(LoadMisc())
}

func TestNetInterfaces(t *testing.T) {
	printInfo(NetInterfaces())
}

func TestNetConnections(t *testing.T) {
	printInfo(NetConnections())
}

func TestNetIOCounters(t *testing.T) {
	printInfo(NetIOCounters())
}

func TestDiskPartitions(t *testing.T) {
	printInfo(DiskPartitions())
}

func TestDiskUsage(t *testing.T) {
	printInfo(DiskUsage())
}
