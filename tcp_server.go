package redis

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/roshats/redis/internal/respgo"
	"net"
	"strings"
	"time"
)

var connTimeout = time.Minute

func handleConnection(conn net.Conn, r CommandProcessor) {
	defer conn.Close()

	bufReader := bufio.NewReader(conn)

	for {
		conn.SetDeadline(time.Now().Add(connTimeout))

		decoded, err := respgo.Decode(bufReader)
		if err != nil {
			break
		}
		arr, err := convertRESPRequest(decoded)
		if err != nil {
			break
		}

		command, query := arr[0], arr[1:]
		command = strings.ToLower(command)
		result := r.ProcessCommand(command, query)
		conn.Write(result.MarshalRESP())
		if _, ok := result.(*quitResultType); ok {
			break
		}
	}
}

func convertRESPRequest(decoded interface{}) ([]string, error) {
	interfaces, ok := decoded.([]interface{})
	if !ok || len(interfaces) < 1 {
		return nil, errors.New("invalid RESP command format")
	}

	result := make([]string, len(interfaces))
	for i := range interfaces {
		result[i], ok = interfaces[i].(string)
		if !ok {
			return nil, errors.New("invalid RESP command format")
		}
	}
	return result, nil
}

func StartTCPServer(s *Server, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		listener.Close()
		fmt.Println("Listener closed")
	}()

	for {
		// Get net.TCPConn object
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}

		client := NewEmbeddedClient(s)
		go handleConnection(conn, client)
	}
}
