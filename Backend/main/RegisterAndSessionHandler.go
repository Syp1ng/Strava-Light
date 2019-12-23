package main

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type sessionKeyInfo struct {
	validUntil time.Time
	forUser    int
}

//Settings
var sessionKeyExpires time.Duration = 24 * 7
var sessionKeyLen int = 50

var allSessions map[string]sessionKeyInfo

func checkSessionKey(sessionKey string) bool {
	fmt.Println(sessionKey)
	info, exists := allSessions[sessionKey]
	if exists {
		delta := info.validUntil.Sub(time.Now())
		fmt.Println(info)
		if delta > 0 {
			fmt.Print("valid")
			return true
		}
	}
	fmt.Println("sesssion not valid")
	return false
}
func generateSessionKey(userID int) string {
	sessionKey := getRandomString(sessionKeyLen)
	allSessions[sessionKey] = sessionKeyInfo{
		validUntil: time.Now().Add(time.Hour * sessionKeyExpires),
		forUser:    userID,
	}
	return base64.StdEncoding.EncodeToString([]byte(sessionKey))
}

func getRandomString(keyLen int) string {
	var charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var generatedKey string
	for i := 0; i <= keyLen; i++ {
		generatedKey += string(charset[rand.Intn(len(charset))])
	}
	return generatedKey
}

func login(username string, password string) string {
	/*sum := sha512.Sum512([]byte(password))
	var passInDB string
	uID:=0
	userData, error := os.Open("userData.csv")
	if error == nil {
		reader := csv.NewReader(userData)
		for {
			line, err := reader.Read()
			if err == nil {
				if line[1] == username{
					passInDB = line[2]
					uID = line[0]
					break
				}
			} else{
				break
			}
		}

	}*/
	return generateSessionKey(0)
}

func register(username string, password string) string {
	f, err := os.Open("UserDataDB.csv")
	if err != nil {
		log.Println(err)
	}

	userData := csv.NewReader(f)
	for {
		line, err := userData.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if line[0] == username {
			log.Println("Username bereits vorhanden")
			return "Username already used"
		}
	}
	//hashing
	//sum := sha512.Sum512([]byte(password))

	writeData := csv.NewWriter(f)

	writeData.Flush()

	return "Succesfull"
}
