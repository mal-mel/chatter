package main

import "time"


const (
	Timeout  = 600 * time.Second
	BuffSize = 4096
	HelpDescription = "id - display your id\n" +
		              "chat <id> - start chat with user with given id\n" +
		              "online - get online user list\n"
)
