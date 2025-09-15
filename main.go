package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Starting Pithos!")
	fmt.Printf("Parent PID: %d\n", os.Getpid())
	fmt.Printf("Parent UID: %d, GID: %d\n", os.Getuid(), os.Getgid())

	//we place the parent pid in the wanted cgroup so the child can immediately inherit it
	cgroupID := "test123" //TODO Config

	if err := setupCgroups(cgroupID, os.Getpid()); err != nil {
		fmt.Printf("Failed to setup cgroups: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := cleanupCgroup(cgroupID); err != nil { //the kernel should take care of this normally, we can prob delete this defer
			fmt.Printf("Warning: failed to cleanup cgroup: %v\n", err)
		}

	}()

	//name space setup
	cmd := exec.Command("/proc/self/exe", "child")
	setupNamespace(cmd)

	// Set up stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting container...")

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start container: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Container error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Container finished successfully!")
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	if len(os.Args) > 1 && os.Args[1] == "child" {
		childFunction()
		os.Exit(0)
	}
}
