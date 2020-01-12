//////////5422223//////////9872387//////////8190324//////////
package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	//allSessions = make(map[string]sessionKeyInfo)
	//userDataMap = make(map[int]userData)
	activityMap = make(map[int]Activity)
	activityMapForUser = make(map[int]Activity)
	/*
		port = ":" + strconv.Itoa(*flag.Int("port", 443, "Port for Webserver"))
		saltLen = *flag.Int("saltLen", 10, "Length of the salt for the Password")
		sessionKeyLen = *flag.Int("sessionKeyLen", 50, "Length of the session Key")
	*/

	//this way because of defining settings in main
	flagPort := flag.Int("port", 443, "Port for Webserver")
	flagSaltLen := flag.Int("saltLen", 10, "Length of the salt for the Password")
	flagSessionKeyLen := flag.Int("sessionKeyLen", 50, "Length of the session Key")

	flag.Parse()

	port = ":" + strconv.Itoa(*flagPort)
	saltLen = *flagSaltLen
	sessionKeyLen = *flagSessionKeyLen

	fmt.Println("Webserver Starting on port" + port)

	SetupLinks()
}
