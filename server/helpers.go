package main


func GetFreeId() int {
	maxId := 0
	for host, _ := range ConnectionsInterfacesAddr {
		if id := ConnectionsInterfacesAddr[host].id; id > maxId {
			maxId = id
		}
	}
	return maxId + 1
}

func IsExistString(collection []string, elem string) bool {
	for _, v := range collection {
		if v == elem {
			return true
		}
	}
	return false
}
