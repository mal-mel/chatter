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
	scanner := bufio.NewScanner(reader)
	for {
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		bData := scanner.Bytes()
		structType := GetStructType(bData)
		if structType == "response" {
			var responseStruct Response
			err := json.Unmarshal(bData, &responseStruct)
			if err != nil {
				fmt.Println(err)
			}
			responseHandler(&responseStruct, connInterface)
		} else if structType == "command" {
			var commandStruct Command
			err := json.Unmarshal(bData, &commandStruct)
			if err != nil {
				fmt.Println(err)
			}
			commandHandler(&commandStruct, connInterface)
		}
	}
}

func responseHandler(responseStruct *Response, connInterface *Conn) {
	if len(responseStruct.Data) > 0 {
		fmt.Println("Command from " + connInterface.GetRemoteAddr() + ": " + responseStruct.Data)
		splittedCommand := strings.Split(responseStruct.Data, " ")
		clientCommand := splittedCommand[0]
		if handler, ok := CommandInterfaces[clientCommand]; ok {
			response := handler(connInterface, splittedCommand[1:])
			err := connInterface.Write(response)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err := connInterface.Write(GetErrorResponse())
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func commandHandler(commandStruct *Command, connInterface *Conn) {
	if commandStruct.Code == 200 {
		RemoveId(commandStruct.Id)
	} else {
		bJson := MakeJson(&commandStruct)
		err := connInterface.Write(bJson)
		if err != nil {
			fmt.Println(err)
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
	id := GetFreeUserId()
	clientAddr := connInterface.GetRemoteAddr()
	ConnectionsInterfacesAddr[clientAddr] = &InterfaceData{id, connInterface}
	ConnectionsInterfacesId[id] = &InterfaceData{id, connInterface}
	fmt.Println("Update interfaces: ", CommandInterfaces)
	return connInterface
}

func main() {
	listen, err := net.Listen("tcp", Address + Port)
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
		err = connInterface.SetDeadline(time.Now().Add(connInterface.IdleTimeout))
		if err != nil {
			fmt.Println(err)
		}
		go requestHandle(connInterface)
	}
}
