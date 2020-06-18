package main

import (
	"encoding/json"
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
	response, err := json.Marshal(&Message{"response", 400, "Error"})
	if err != nil {
		fmt.Println(err)
	}
	response = append(response, "\n"...)
	return response
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
	if iData, ok := ConnectionsInterfacesAddr[connInterface.GetRemoteAddr()]; ok {
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
	onlineIds := getOnlineIds(connInterface)
	bJson, err := json.Marshal(&Message{"response", 200, strings.Join(onlineIds[:], "\n")})
	if err != nil {
		panic(err)
	}
	bJson = append(bJson, "\n"...)
	return bJson
}

func getOnlineIds(connInterface *Conn) []string {
	var result []string
	var clientId int
	if clientData, ok := ConnectionsInterfacesAddr[connInterface.GetRemoteAddr()]; ok {
		clientId = clientData.id
	}
	for iData := range ConnectionsInterfacesAddr {
		if ConnectionsInterfacesAddr[iData].id != clientId {
			result = append(result, strconv.Itoa(ConnectionsInterfacesAddr[iData].id))
		}
	}
	return result
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
					bJson, err := json.Marshal(&Message{CurrentId, "response", 200, "accept_" + strconv.Itoa(c1Data.id)})
					CurrentId++
					if err != nil {
						fmt.Println(err)
					}
					bJson = append(bJson, "\n"...)
					_, err = c2Data.conn.Write(bJson)
					if err != nil {
						fmt.Println(err)
					}
					var buff []byte
					var responseC2 Message
					_, err = c2Data.conn.Read(buff)
					if err != nil {
						fmt.Println(err)
					}
					err = json.Unmarshal(buff, &responseC2)
					buff = nil
					if responseC2.Code == 200 {
						return buff
					}
				}
			}
		}
	}
	return GetErrorResponse()
}
