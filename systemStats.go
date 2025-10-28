package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
    logicalCores, _ := cpu.Counts(true)
	physicalCores, _ := cpu.Counts(false)
	info, _ := cpu.Info()
	fmt.Printf("%s\n", info[0].ModelName)

	fmt.Printf("%v Physical Cores\n", physicalCores)
	fmt.Printf("%v Threads\n", logicalCores)

    platform, _ ,version, _:= host.PlatformInformation()
    fmt.Println("Platform:", platform)
    fmt.Println("Version:", version)

	m, _ := mem.VirtualMemory()
	totalMem := m.Total / (1024 * 1024 * 1024)
	fmt.Printf("%vGb Total Usable RAM\n", totalMem)

	d, _ := disk.Partitions(false)
	usage, _ := disk.Usage(d[0].Mountpoint)
	freeDisk := usage.Free / (1024 * 1024 * 1024)
	totalDisk := usage.Total / (1024 * 1024 * 1024)
	fmt.Printf("%vGb Total Disc Space\n", totalDisk)
	fmt.Printf("%vGb Free Disc Space\n", freeDisk)
	
}
