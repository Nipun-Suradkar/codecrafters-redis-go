package replication

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/redis_server"
)

func InformMasterServer(redisServer *redis_server.RedisServer) {
	replicaOf := redisServer.ReplicaOf
	if replicaOf == "" {
		return
	}

	masterInfo := strings.Split(replicaOf, " ")
	if len(masterInfo) != 2 {
		return
	}

	masterHost := masterInfo[0]
	masterPort := masterInfo[1]
	masterAddr := fmt.Sprintf("%s:%s", masterHost, masterPort)

	conn, err := net.Dial("tcp", masterAddr) // host:port
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	writer := bufio.NewWriter(conn)
	writeArrayBulk(writer, "PING")
	defer conn.Close()

}

func writeArrayBulk(writer *bufio.Writer, msgs ...string) {
	writer.WriteString(fmt.Sprintf("*%d\r\n", len(msgs)))
	for _, s := range msgs {
		writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	}
	writer.Flush()
}
