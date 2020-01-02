package main

import (
	"fmt"
)

func main() {
	fmt.Println("Webserver Starting")
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	SetupLinks()
	dropTable()
	/* for later
	var port = flag.Int("port", 443, "Port for Webserver")
	var saltLen = flag.Int("saltLen", 1234, "Length of the salt for the Password")
	var sessionKeyLen = flag.Int("sessionKeyLen", 50, "Length of the session Key")
	var sessionKeyExpires = flag.Int("sessionKeyExpires", 24*7, "Time in houres how long session Key is valid")
	*/

}
