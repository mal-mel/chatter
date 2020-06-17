package main

import (
	"bufio"
	"encoding/json"
	"fmt"
)


type handlersType func(*bufio.Scanner, *bufio.Writer, []byte)


var CommandInterfaces = map[string]handlersType {
	"chat": chatCommandHandler,
}


func chatCommandHandler(scanner *bufio.Scanner, serverWriter *bufio.Writer, bJson []byte) {
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
	fmt.Println(messageStruct)
}
