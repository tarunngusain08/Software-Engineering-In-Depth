package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	// Check if bpftrace is installed
	_, err := exec.LookPath("bpftrace")
	if err != nil {
		log.Fatalf("bpftrace is not installed. Please install it using your package manager.")
	}

	// Define the eBPF program
	bpfProgram := `
	tracepoint:syscalls:sys_enter_connect {
		// Extract the sockaddr structure from args->uservaddr
		// For IPv4, the port is stored in the second 2 bytes of the sockaddr_in structure
		$addr = args->uservaddr;
		$port = (uint16)(*(uint16 *)($addr + 2));
		printf("PID: %d, Comm: %s, Port: %d\n", pid, comm, $port);
	}
	`

	// Write the program to a temporary file
	tmpFile, err := os.CreateTemp("", "ebpf-program-*.bt")
	if err != nil {
		log.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(bpfProgram)
	if err != nil {
		log.Fatalf("Failed to write eBPF program to file: %v", err)
	}
	tmpFile.Close()

	// Run the eBPF program using bpftrace
	cmd := exec.Command("bpftrace", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running eBPF program to monitor active connections...")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run eBPF program: %v", err)
	}
}
