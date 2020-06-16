package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)


var ConnectionsInterfaces = make(map[string]*InterfaceData)


func requestHandle(connInterface *Conn) {
	defer func() {
		err := connInterface.Conn.Close()
		if err != nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(connInterface.Conn)
	writer := bufio.NewWriter(connInterface.Conn)
	scanner := bufio.NewScanner(reader)
	for {
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		request := scanner.Bytes()
		var message Message
		err := json.Unmarshal(request, &message)
		if err != nil {
			panic(err)
		}
		if len(message.Data) > 0 {
			fmt.Println("Command from " + connInterface.GetRemoteAddr() + ": " + message.Data)
			splittedCommand := strings.Split(message.Data, " ")
			clientCommand := splittedCommand[0]
			if handler, ok := CommandInterfaces[clientCommand]; ok {
				response := handler(connInterface, splittedCommand[1:])
				_, err = writer.Write(response)
				if err != nil {
					panic(err)
				}
			} else {

			}
			_ = writer.Flush()
		}
	}
}

func initConnectionInterface(mainConn net.Conn) *Conn {
	connInterface := &Conn{
		Conn:          mainConn,
		IdleTimeout:   Timeout,
		MaxReadBuffer: BuffSize,
	}
	_ = connInterface.SetDeadline(time.Now().Add(connInterface.IdleTimeout))
	id := GetFreeId()
	clientAddr := connInterface.GetRemoteAddr()
	ConnectionsInterfaces[clientAddr] = &InterfaceData{id, connInterface}
	fmt.Println("Update interfaces: ", CommandInterfaces)
	return connInterface
}

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Cant start server")
	}
	for {
		client, err := listen.Accept()
		if err != nil {
			fmt.Println("Cant accept connection")
		}
		connInterface := initConnectionInterface(client)
		fmt.Println("Connect from: " + connInterface.GetRemoteAddr())
		connInterface.SetDeadline(time.Now().Add(connInterface.IdleTimeout))
		go requestHandle(connInterface)
	}
}