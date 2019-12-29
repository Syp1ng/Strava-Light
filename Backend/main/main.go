package main

import (
	"fmt"
)

func main() {
	fmt.Println("Webserver Starting")
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	SetupLinks()
}
