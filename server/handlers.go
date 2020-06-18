package main

import (
	"fmt"
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


func GetErrorResponse() [] byte {
	bJson := MakeJson(Response{400, "Error"})
	return bJson
}

func helpCommandHandler(connInterface *Conn, args []string) []byte {
	bJson := MakeJson(Response{200, HelpDescription})
	return bJson
}

func idCommandHandler(connInterface *Conn, args []string) []byte {
	if iData, ok := ConnectionsInterfacesAddr[connInterface.GetRemoteAddr()]; ok {
		bJson := MakeJson(Response{200, strconv.Itoa(iData.id)})
		return bJson
	} else {
		bJson := MakeJson(Response{400, "Unknown user"})
		return bJson
	}
}

func onlineCommandHandler(connInterface *Conn, args []string) []byte {
	onlineIds := GetOnlineIds(connInterface)
	bJson := MakeJson(Response{200, strings.Join(onlineIds[:], "\n")})
	return bJson
}

func chatCommandHandler(connInterface *Conn, args []string) []byte {
	if len(args) > 0 {
		c2Id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
		}
		if c2Data, ok := ConnectionsInterfacesId[c2Id]; ok {
			if c1Data, ok := ConnectionsInterfacesAddr[connInterface.GetRemoteAddr()]; ok {
				if c1Data.id != c2Data.id {
					bJson := MakeJson(Command{GetId(), 200, "ACCEPT_" + strconv.Itoa(c1Data.id)})
					IncrementId()
					err = c2Data.conn.Write(bJson)
					if err != nil {
						fmt.Println(err)
					}
					statusCode, buff := c2Data.GetClientResponse()
					if statusCode == 200 {
						return buff
					}
				}
			}
		}
	}
	return GetErrorResponse()
}
