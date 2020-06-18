package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)


func GetFreeUserId() int {
	maxId := 0
	for host, _ := range ConnectionsInterfacesAddr {
		if id := ConnectionsInterfacesAddr[host].id; id > maxId {
			maxId = id
		}
	}
	return maxId + 1
}


func GetStructType(request []byte) string{
	var response Response; var command Command
	err := json.Unmarshal(request, &response)
	if err == nil {
		return "response"
	}
	err = json.Unmarshal(request, &command)
	if err == nil {
		return "command"
	}
	return ""
}


func GetOnlineIds(connInterface *Conn) []string {
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


func IsExistString(collection []string, elem string) bool {
	for _, v := range collection {
		if v == elem {
			return true
		}
	}
	return false
}


func MakeJson(i interface{}) []byte {
	bJson, err := json.Marshal(&i)
	if err != nil {
		fmt.Println(err)
	}
	bJson = append(bJson, "\n"...)
	return bJson
}
