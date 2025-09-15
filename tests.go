package main

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func childFunction() {
	fmt.Println("Hello from Pithos container!")
	fmt.Printf("Container PID: %d\n", os.Getpid())
	fmt.Printf("Container PPID: %d\n", os.Getppid())
	fmt.Printf("Container UID: %d, GID: %d\n\n", os.Getuid(), os.Getgid())

	//Verify we're in the cgroup by reading /proc/self/cgroup
	fmt.Println("Checking cgroup membership:")
	if cgroupData, err := os.ReadFile("/proc/self/cgroup"); err != nil {
		fmt.Printf("Failed to read /proc/self/cgroup: %v\n", err)
	} else {
		fmt.Printf("Current cgroup: %s", string(cgroupData))
	}

	//Set hostname to show UTS isolation
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf("HostName is: %s before change\n", hostname)

	if err := syscall.Sethostname([]byte("pithos-container")); err != nil {
		fmt.Printf("Failed to set hostname: %v\n", err)
	} else {
		fmt.Printf("Set hostname to: pithos-container\n")
	}

	hostnameNew, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf("HostName is: %s after change\n\n", hostnameNew)

	//Verify IPC isolation: print current message queues (should be empty/new)
	fmt.Println("IPC check: running `ipcs -q`")
	_ = runCommand("ipcs", "-q")

	//Verify NET isolation: show interfaces (should only see loopback down initially)
	fmt.Println("NET check: running `ip addr`")
	_ = runCommand("ip", "addr")

	//Test filesystem access
	if _, err := os.ReadFile("/etc/passwd"); err != nil {
		fmt.Printf("\nCHECK Can't read /etc/passwd: %v\n", err)
	} else {
		fmt.Printf("\nWARNING Can still read /etc/passwd\n")
	}

	fmt.Println("Container sleeping for 3 seconds...")
	time.Sleep(3 * time.Second)
	fmt.Println("Container finished!")
}
