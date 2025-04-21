package ebpf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/perf"
)

func main() {
	// Load the compiled eBPF program from the .o file
	spec, err := ebpf.LoadCollectionSpec("monitor_incoming.o")
	if err != nil {
		log.Fatalf("Failed to load eBPF collection spec: %v", err)
	}

	// Load the program into the kernel
	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		log.Fatalf("Failed to create eBPF collection: %v", err)
	}
	defer coll.Close()

	prog := coll.Programs["monitor_incoming"]
	if prog == nil {
		log.Fatalf("Program 'monitor_incoming' not found in the collection")
	}
	defer prog.Close()

	// Get the events map
	eventsMap, ok := coll.Maps["events"]
	if !ok {
		log.Fatalf("Map 'events' not found in the collection")
	}

	// Debug: Ensure eventsMap is not nil
	if eventsMap == nil {
		log.Fatalf("eventsMap is nil")
	}

	// // Create a perf buffer to receive events
	reader, err := perf.NewReader(eventsMap, 64) // 64 pages for the buffer
	if err != nil {
		log.Fatalf("Failed to create perf reader: %v", err)
	}
	defer reader.Close()

	go func() {
		for {
			record, err := reader.Read()
			if err != nil {
				log.Printf("Failed to read from perf buffer: %v", err)
				continue
			}

			var conn struct {
				PID   uint32
				SAddr uint32
				DAddr uint32
				SPort uint16
				DPort uint16
			}
			if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &conn); err != nil {
				log.Printf("Failed to decode perf event: %v", err)
				continue
			}

			fmt.Printf("PID: %d, Source: %s:%d, Destination: %s:%d\n",
				conn.PID,
				intToIP(conn.SAddr), conn.SPort,
				intToIP(conn.DAddr), conn.DPort)
		}
	}()

	// Handle termination signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func intToIP(ip uint32) net.IP {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}
