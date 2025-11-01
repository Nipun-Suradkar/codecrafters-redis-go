package redis_server

import (
	"flag"
	"log"
	"strconv"
)

var serverConfig *RedisServer

type RedisServer struct {
	DbFilename    string
	DbDir         string
	Port          int
	ReplicaOf     string
	ReplicationID string
	Offset        int32
	IsSlaveNode   bool
}

func GetRedisServer() *RedisServer {
	return serverConfig
}

func InitializeRedisServer() {
	dir := flag.String("dir", "", "Directory path for Redis data")
	dbFilename := flag.String("dbfilename", "", "Database filename")
	portString := flag.String("port", "6379", "Port to accept connections")
	replicaOf := flag.String("replicaof", "", "Address of master (for replica mode)")
	flag.Parse()

	port, err := strconv.Atoi(*portString)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	isSalveNode := *replicaOf != ""

	serverConfig = &RedisServer{
		DbDir:         *dir,
		DbFilename:    *dbFilename,
		ReplicaOf:     *replicaOf,
		Port:          port,
		ReplicationID: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb", // fixed unique value
		Offset:        0,
		IsSlaveNode:   isSalveNode,
	}
}
