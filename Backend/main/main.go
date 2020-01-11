package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("Webserver Starting")
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	activityMap = make(map[int]Activity)
	activityMapForUser = make(map[int]Activity)

	port = ":" + strconv.Itoa(*flag.Int("port", 443, "Port for Webserver"))
	saltLen = *flag.Int("saltLen", 10, "Length of the salt for the Password")
	sessionKeyLen = *flag.Int("sessionKeyLen", 50, "Length of the session Key")

	SetupLinks()
}
