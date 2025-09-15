package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func setupCgroups(containerID string, pid int) error {
	cgroupPath := filepath.Join("/sys/fs/cgroup", "pithos-container-"+containerID)

	//Create cgroup directory
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup: %v", err)
	}
	fmt.Printf("Created cgroup: %s\n", cgroupPath)

	//Set memory limit to 128MB
	memoryMax := filepath.Join(cgroupPath, "memory.max")
	if err := os.WriteFile(memoryMax, []byte("134217728"), 0644); err != nil { //TODO config
		return fmt.Errorf("failed to set memory limit: %v", err)
	}
	fmt.Println("Set memory limit: 128MB")

	//Set CPU weight
	cpuWeight := filepath.Join(cgroupPath, "cpu.weight")
	if err := os.WriteFile(cpuWeight, []byte("100"), 0644); err != nil { //TODO config
		return fmt.Errorf("failed to set CPU weight: %v", err)
	}
	fmt.Println("Set CPU weight: 100")

	pidsMax := filepath.Join(cgroupPath, "pids.max")
	if err := os.WriteFile(pidsMax, []byte("128"), 0644); err != nil { //TODO config
		return fmt.Errorf("failed to set pids.max: %v", err)
	}
	fmt.Println("Set pids.max: 128")

	//Add process to cgroup
	cgroupProcs := filepath.Join(cgroupPath, "cgroup.procs")
	if err := os.WriteFile(cgroupProcs, []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %v", err)
	}
	fmt.Printf("Added PID %d to cgroup\n", pid)

	return nil
}

func cleanupCgroup(containerID string) error {
	cgroupPath := filepath.Join("/sys/fs/cgroup", "pithos-container-"+containerID)

	if err := os.RemoveAll(cgroupPath); err != nil {
		return fmt.Errorf("failed to remove cgroup: %v", err)
	}

	fmt.Printf("Cleaned up cgroup: %s\n", cgroupPath)
	return nil
}
