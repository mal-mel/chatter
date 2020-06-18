package main

import "time"


const (
	Address = "localhost"
	Port = ":8080"

	Timeout  = 600 * time.Second
	BuffSize = 4096

	HelpDescription = "id - display your id\n" +
		              "chat <id> - start chat with user with given id\n" +
		              "online - get online user list\n"
)
