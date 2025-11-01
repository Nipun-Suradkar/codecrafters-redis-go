package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/redis_server"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	log.Println("Server listening on port 6379...")

	redis_server.InitializeRedisServer()

	if redis_server.GetRedisServer().IsSlaveNode {
		replication.InformMasterServer()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", redis_server.GetRedisServer().Port))
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	commandHandler := command.NewCommand(writer)

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
