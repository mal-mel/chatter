package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)


var ConnectionsInterfacesAddr = make(map[string]*InterfaceData)
var ConnectionsInterfacesId = make(map[int]*InterfaceData)


func requestHandle(connInterface *Conn) {
	defer func() {
		err := connInterface.Conn.Close()
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
		}
		if len(message.Data) > 0 {
			fmt.Println("Command from " + connInterface.GetRemoteAddr() + ": " + message.Data)
			splittedCommand := strings.Split(message.Data, " ")
			clientCommand := splittedCommand[0]
			if handler, ok := CommandInterfaces[clientCommand]; ok {
				response := handler(connInterface, splittedCommand[1:])
				_, err = writer.Write(response)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				errorResponse, err := GetErrorResponse()
				if err != nil {
					fmt.Println(err)
				}
				_, err = writer.Write(errorResponse)
				if err != nil {
					fmt.Println(err)
				}
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
	ConnectionsInterfacesAddr[clientAddr] = &InterfaceData{id, connInterface}
	ConnectionsInterfacesId[id] = &InterfaceData{id, connInterface}
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
