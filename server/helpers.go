package main


func GetFreeId() int {
	maxId := 0
	for host, _ := range ConnectionsInterfaces {
		if id := ConnectionsInterfaces[host].id; id > maxId {
			maxId = id
		}
	}
	return maxId + 1
}
