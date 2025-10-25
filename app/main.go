package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/data_store"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	log.Println("Server listening on port 6379...")

	dir := flag.String("dir", "", "Directory path")
	dbFilename := flag.String("dbfilename", "", "Database filename")
	flag.Parse()
	dataStore := &data_store.DataStore{
		DbFilename: *dbFilename,
		DbDir:      *dir,
	}

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to bind to port 6379: %v", err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Println("Failed to close listener:", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, dataStore)
	}
}

func handleConnection(conn net.Conn, store *data_store.DataStore) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	commandHandler := command.NewCommand(writer, store)

	for {
		// Try to decode a RESP array from the reader
		cmd, err := resp.DecodeRESPFromReader(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
			} else {
				fmt.Println("Error decoding RESP:", err)
			}
			return
		}
		fmt.Println("Decoded command:", cmd)

		commandHandler.HandleCommand(cmd)

	}
}
