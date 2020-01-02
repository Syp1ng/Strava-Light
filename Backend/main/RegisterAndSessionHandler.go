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
var sessionKeyExpires time.Duration = 24 * 7         //how long the session Key is valid
var sessionKeyLen int = 50                           //length of the session Key
var saltLen int = 10                                 //Length of the salt for the Password
var dbLocation string = "DataStorage/UserDataDB.csv" //Path to the UserData

var allSessions map[string]sessionKeyInfo
var userDataMap map[int]userData

type userData struct {
	uID      int
	userName string
	email    string
	password string
}

//Function to get the userID from the sessionKey | checks vor valid SessionKey is done before calling that function
func getUID(sessionKey string) int {
	return allSessions[sessionKey].forUser
}

//checks if a sessionkey is valid and not expired
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

//function which returns a sessionkey and saves the key on server side in a map, with more details
func generateSessionKey(userID int) string {
	sessionKey := getRandomString(sessionKeyLen)
	expiresOn := time.Now().Add(time.Hour * sessionKeyExpires)
	allSessions[sessionKey] = sessionKeyInfo{
		validUntil: expiresOn,
		forUser:    userID,
	}
	return sessionKey //base64.StdEncoding.EncodeToString([]byte(sessionKey))
}

//deletes a sessionkey / makes it invalid
func delSessionKey(sessionKey string) {
	delete(allSessions, sessionKey)
}

//Function which returns a String of a costum length from a charset
func getRandomString(keyLen int) string {
	rand.Seed(time.Now().UnixNano())
	var charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var generatedKey string
	for i := 0; i <= keyLen; i++ {
		generatedKey += string(charset[rand.Intn(len(charset))])
	}
	return generatedKey
}

//Function to Login the user and return 1) if it was successfull 2) an error Message
func login(userName string, password string) (bool, string) {
	userData, error := os.Open(dbLocation)
	if error == nil {
		reader := csv.NewReader(userData)
		for {
			line, err := reader.Read()
			if err == nil {
				if line[1] == userName {
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

//function to compare the cleartext Password from User with the hashed Password form the DB
func comparePasswords(userInputPass, dBPassAndSalt string) bool {
	hashAlgo := sha512.New()
	dBPassAndSaltArray := strings.Split(dBPassAndSalt, ":")
	hashAlgo.Write([]byte(userInputPass + dBPassAndSaltArray[1]))
	if hex.EncodeToString(hashAlgo.Sum(nil)) == dBPassAndSaltArray[0] {
		return true
	}
	return false
}

//takes the password and adds a salt, and returns the hash of this with the salt after the :
func hashPassword(password string) string {
	hashAlgo := sha512.New()
	salt := getRandomString(saltLen)
	hashAlgo.Write([]byte(password + salt))
	return hex.EncodeToString(hashAlgo.Sum(nil)) + ":" + salt
}

func register(userName string, email string, password string, confirmPass string) (bool, string) {
	if password != confirmPass {
		return false, "Passwörter stimmen nicht überein"
	} else if len(password) < 8 {
		return false, "Bitte ein längeres Passwort"
	}

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
	newUser := userData{maxID + 1, userName, email, hashPassword(password)}
	userDataMap[maxID+1] = newUser
	if appendToDB(newUser) {
		return true, generateSessionKey(newUser.uID)
	} else {
		return false, "error writing to db"
	}

}

func appendToDB(user userData) bool {
	var newline string = strconv.Itoa(user.uID) + "," + user.userName + "," + user.email + "," + user.password + "\n"
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

func dropTable() {
	err := os.Remove(dbLocation)
	if err != nil {
		fmt.Println("cannot detele file")
	}
	_, err2 := os.Create(dbLocation)
	if err2 != nil {
		fmt.Println("cannot create file")
	}

}
