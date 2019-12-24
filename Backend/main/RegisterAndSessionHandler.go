package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type sessionKeyInfo struct {
	validUntil time.Time
	forUser    int
}

//Settings
var sessionKeyExpires time.Duration = 24 * 7
var sessionKeyLen int = 50
var saltLen int = 10
var dbLocation string = "Backend/main/UserDataDB.csv"

var allSessions map[string]sessionKeyInfo
var userDataMap map[int]userData

type userData struct {
	uID      int
	email    string
	name     string
	password string
}

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
	expiresOn := time.Now().Add(time.Hour * sessionKeyExpires)
	allSessions[sessionKey] = sessionKeyInfo{
		validUntil: expiresOn,
		forUser:    userID,
	}
	return sessionKey //base64.StdEncoding.EncodeToString([]byte(sessionKey))
}

func getRandomString(keyLen int) string {
	var charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var generatedKey string
	for i := 0; i <= keyLen; i++ {
		generatedKey += string(charset[rand.Intn(len(charset))])
	}
	return generatedKey
}

func login(email string, password string) (bool, string) {
	userData, error := os.Open(dbLocation)
	if error == nil {
		reader := csv.NewReader(userData)
		for {
			line, err := reader.Read()
			if err == nil {
				if line[1] == email {
					if comparePasswords(password, line[3]) {
						x, _ := strconv.Atoi(line[0])
						return true, generateSessionKey(x)
					} else {
						return false, "Wrong password"
					}
				}
			} else {
				break
			}
		}
	}
	return false, "Email unknown"
}

func comparePasswords(userInputPass, dBPassAndSalt string) bool {
	hashAlgo := sha512.New()
	dBPassAndSaltArray := strings.Split(dBPassAndSalt, ":")
	hashAlgo.Write([]byte(userInputPass + dBPassAndSaltArray[1]))
	if hex.EncodeToString(hashAlgo.Sum(nil)) == dBPassAndSaltArray[0] {
		return true
	}
	return false
}
func hashPassword(password string) string {
	hashAlgo := sha512.New()
	salt := getRandomString(saltLen)
	hashAlgo.Write([]byte(password + salt))
	return hex.EncodeToString(hashAlgo.Sum(nil)) + ":" + salt
}

func register(email string, username string, password string) (bool, string) {
	readDB()
	var maxID int = 0
	for k, v := range userDataMap {
		if email == v.email {
			return false, "Account exist with that email"
		}
		if k > maxID {
			maxID = k
		}

	}
	newUser := userData{maxID + 1, email, username, hashPassword(password)}
	userDataMap[maxID+1] = newUser
	if appendToDB(newUser) {
		return true, generateSessionKey(newUser.uID)
	} else {
		return false, "error writing to db"
	}

}

func appendToDB(user userData) bool {
	var newline string = strconv.Itoa(user.uID) + "," + user.email + "," + user.name + "," + user.password + "\n"
	f, err := os.OpenFile(dbLocation, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.WriteString(newline)
	if err != nil {
		return false
	}
	return true
}
func readDB() {
	file, err := os.Open(dbLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		user := strings.Split(scanner.Text(), ",")
		if len(user) != 4 {
			continue //invalid or  empty
		}
		userID, _ := strconv.Atoi(user[0])
		newUser := userData{userID, user[1], user[2], user[3]}
		userDataMap[userID] = newUser
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
