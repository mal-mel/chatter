package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)


func makeMessage(data string) []byte {
	data = data[:len(data) - 1]
	message := Message{"command", 0, data}
	bJson, err := json.Marshal(&message)
	if err != nil {
		fmt.Println(err)
	}
	bJson = append(bJson, "\n"...)
	return bJson
}

func getResponse(scanner *bufio.Scanner, serverWriter *bufio.Writer, bJson []byte) Message {
	var messageStruct Message
	_, err := serverWriter.Write(bJson)
	if err != nil {
		fmt.Println(err)
	}
	err = serverWriter.Flush()
	if err != nil {
		fmt.Println(err)
	}
	_ = scanner.Scan()
	err = json.Unmarshal(scanner.Bytes(), &messageStruct)
	if err != nil {
		fmt.Println(err)
	}
	return messageStruct
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
	}
	serverReader := bufio.NewReader(conn)
	serverWriter := bufio.NewWriter(conn)
	scanner := bufio.NewScanner(serverReader)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("chatter -> ")
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if len(data) > 1 {
			bJson := makeMessage(data)
			response := getResponse(scanner, serverWriter, bJson)
			if len(response.Data) != 0 {
				fmt.Println(response.Data)
			}
		}
	}
}
