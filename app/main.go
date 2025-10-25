package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/redis_server"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	log.Println("Server listening on port 6379...")

	redisServer := initializeRedisServer()

	if redisServer.ReplicaOf != "" {
		replication.InformMasterServer(redisServer)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", redisServer.Port))
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
		go handleConnection(conn, redisServer)
	}
}

func initializeRedisServer() *redis_server.RedisServer {
	dir := flag.String("dir", "", "Directory path")
	dbFilename := flag.String("dbfilename", "", "Database filename")
	portString := flag.String("port", "", "port to accept conections")
	replicaOf := flag.String("replicaof", "", "replica of server")
	flag.Parse()

	redisServer := &redis_server.RedisServer{
		DbFilename:    *dbFilename,
		DbDir:         *dir,
		ReplicaOf:     *replicaOf,
		ReplicationID: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
	}

	portNum := 6379 //default port
	if portString != nil && *portString != "" {
		port, err := strconv.Atoi(*portString)
		if err != nil {
			log.Fatalf("Invalid port number: %v", err)
		}
		portNum = port
	}
	redisServer.Port = portNum
	return redisServer
}

func handleConnection(conn net.Conn, store *redis_server.RedisServer) {
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
