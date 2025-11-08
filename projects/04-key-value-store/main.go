package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alyxpink/go-training/kvstore/persistence"
	"github.com/alyxpink/go-training/kvstore/protocol"
	"github.com/alyxpink/go-training/kvstore/store"
)

var (
	port        = flag.Int("port", 6380, "Server port")
	dataDir     = flag.String("data-dir", "./data", "Data directory")
	snapshotInt = flag.Duration("snapshot-interval", 5*time.Minute, "Snapshot interval")
)

func main() {
	flag.Parse()

	// TODO: Initialize store
	kvStore := store.NewKVStore()

	// TODO: Initialize persistence
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatal(err)
	}

	wal, err := persistence.NewWAL(*dataDir + "/wal.log")
	if err != nil {
		log.Fatal(err)
	}
	defer wal.Close()

	// TODO: Recover from WAL/snapshot
	snapshot := persistence.NewSnapshotManager(*dataDir, *snapshotInt)
	if err := snapshot.LoadLatest(kvStore); err != nil {
		log.Printf("No snapshot found, starting fresh: %v", err)
	}

	if err := wal.Replay(kvStore); err != nil {
		log.Fatal(err)
	}

	// TODO: Start snapshot routine
	go snapshot.Run(kvStore)

	// TODO: Start TCP server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Server listening on :%d", *port)

	// TODO: Handle shutdown gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		snapshot.CreateSnapshot(kvStore)
		os.Exit(0)
	}()

	// TODO: Accept connections
	handler := protocol.NewHandler(kvStore, wal)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn, handler)
	}
}

func handleConnection(conn net.Conn, handler *protocol.Handler) {
	defer conn.Close()
	log.Printf("New connection from %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Text()
		response := handler.Handle(line)
		writer.WriteString(response + "\r\n")
		writer.Flush()
	}
}
