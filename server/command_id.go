package main


var CurrentId = 1
var Ids = make(map[int]int)


func GetId() int {
	Ids[CurrentId] = CurrentId
	return CurrentId
}

func IncrementId() {
	CurrentId++
}

func RemoveId(id int) {
	if _, ok := Ids[id]; ok {
		delete(Ids, id)
	}
}
