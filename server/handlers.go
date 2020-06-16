package main

import (
	"encoding/json"
	"strconv"
	"strings"
)


type handlersType func(*Conn, []string) []byte


var CommandInterfaces = map[string]handlersType {
	"help": helpCommandHandler,
	"id": idCommandHandler,
	"online": onlineCommandHandler,
	"chat": chatCommandHandler,
}


func helpCommandHandler(connInterface *Conn, args []string) []byte {
	bJson, err := json.Marshal(&Message{"response", 200, HelpDescription})
	if err != nil {
		panic(err)
	}
	bJson = append(bJson, "\n"...)
	return bJson
}

func idCommandHandler(connInterface *Conn, args []string) []byte {
	if iData, ok := ConnectionsInterfaces[connInterface.GetRemoteAddr()]; ok {
		bJson, err := json.Marshal(&Message{"response", 200, strconv.Itoa(iData.id)})
		if err != nil {
			panic(err)
		}
		bJson = append(bJson, "\n"...)
		return bJson
	} else {
		bJson, err := json.Marshal(&Message{"response", 400, "Unknown user"})
		if err != nil {
			panic(err)
		}
		bJson = append(bJson, "\n"...)
		return bJson
	}
}

func onlineCommandHandler(connInterface *Conn, args []string) []byte {
	var result []string
	var clientId int
	if clientData, ok := ConnectionsInterfaces[connInterface.GetRemoteAddr()]; ok {
		clientId = clientData.id
	}
	for iData := range ConnectionsInterfaces {
		if ConnectionsInterfaces[iData].id != clientId {
			result = append(result, strconv.Itoa(ConnectionsInterfaces[iData].id))
		}
	}
	bJson, err := json.Marshal(&Message{"response", 200, strings.Join(result[:], "\n")})
	if err != nil {
		panic(err)
	}
	bJson = append(bJson, "\n"...)
	return bJson
}

func chatCommandHandler(connInterface *Conn, args []string) []byte {
	bJson, err := json.Marshal(&Message{"response", 200, "test"})
	if err != nil {
		panic(err)
	}
	bJson = append(bJson, "\n"...)
	return bJson
}
