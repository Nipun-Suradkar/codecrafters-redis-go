package replication

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/redis_server"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func InformMasterServer() {
	replicaOf := redis_server.GetRedisServer().ReplicaOf
	if replicaOf == "" {
		return
	}

	masterInfo := strings.Split(replicaOf, " ")
	if len(masterInfo) != 2 {
		return
	}

	log.Println("Informing master server")

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

	sendSyncMsgsToMaster(writer, reader)
	defer conn.Close()

}

func sendSyncMsgsToMaster(writer *bufio.Writer, reader *bufio.Reader) error {
	syncCmds := [][]string{
		{"PING"},
		{"REPLCONF", "listening-port", strconv.Itoa(redis_server.GetRedisServer().Port)},
		{"REPLCONF", "capa", "psync2"},
		{"PSYNC", "?", "-1"},
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
	resp.WriteArrayBulk(writer, cmd...)
	writer.Flush()
	return parseMasterResponse(reader)
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
