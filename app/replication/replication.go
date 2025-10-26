package replication

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
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
	reader := bufio.NewReader(conn)

	sendSyncMsgsToMaster(writer, reader, redisServer)
	defer conn.Close()

}

func sendSyncMsgsToMaster(writer *bufio.Writer, reader *bufio.Reader, server *redis_server.RedisServer) error {
	syncCmds := [][]string{
		{"PING"},
		{"REPLCONF", "listening-port", strconv.Itoa(server.Port)},
		{"REPLCONF", "capa", "psync2"},
	}

	for _, cmd := range syncCmds {
		if err := sendAndValidate(writer, reader, cmd...); err != nil {
			return fmt.Errorf("sync step failed for %v: %w", cmd, err)
		}
	}

	return nil
}

// sendAndValidate writes a command and checks for valid master response
func sendAndValidate(writer *bufio.Writer, reader *bufio.Reader, cmd ...string) error {
	if err := writeArrayBulk(writer, cmd...); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}
	return parseMasterResponse(reader)
}

// writeArrayBulk encodes and sends a Redis RESP array
func writeArrayBulk(writer *bufio.Writer, parts ...string) error {
	if _, err := fmt.Fprintf(writer, "*%d\r\n", len(parts)); err != nil {
		return err
	}
	for _, part := range parts {
		if _, err := fmt.Fprintf(writer, "$%d\r\n%s\r\n", len(part), part); err != nil {
			return err
		}
	}
	return writer.Flush()
}

// parseMasterResponse reads and validates a Redis master response
func parseMasterResponse(reader *bufio.Reader) error {
	resp, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	resp = strings.TrimSpace(resp)

	switch {
	case strings.HasPrefix(resp, "+OK"),
		strings.HasPrefix(resp, "+PONG"),
		strings.HasPrefix(resp, "+FULLRESYNC"),
		strings.HasPrefix(resp, "+CONTINUE"):
		log.Printf("[Master Response] %s", resp)
		return nil
	default:
		return fmt.Errorf("unexpected master response: %q", resp)
	}
}
