//////////5422223//////////9872387//////////8190324//////////
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

type sessionKeyInfo struct { //simple sessionKeyInfo which is the value of the sessionkey in the map
	validUntil time.Time
	forUser    int
}

func init() {
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	rand.Seed(time.Now().UTC().UnixNano()) // for more randomness
}

//Settings
var sessionKeyLen int //length of the session Key
var saltLen int       //Length of the salt for the Password

var sessionKeyExpires time.Duration = 24 * 7         //how long the session Key is valid
var dbLocation string = "DataStorage/UserDataDB.csv" //Path to the UserData

var allSessions map[string]sessionKeyInfo //saves sessions @ runtime
var userDataMap map[int]userData          //csv -> map ->csv

//Error Messages for User
var ErrorMessageRegisterNotSamePass = "Fehler: Passwörter stimmen nicht überein"
var ErrorMessageRegisterInvalidPasswordPolicy = "Fehler: Ausgewähltes Passwort zu klein"
var ErrorMessageRegisterUsernameTaken = "Fehler: Benutzername schon vergeben"
var ErrorMessageLoginUsernameUnknown = "Fehler: Benutzername existiert nicht"
var ErrorMessageLoginPasswordWrong = "Fehler: Passwort ist falsch"
var ErrorUnknown = "Datenbankfehler oder Programmfehler"

type userData struct { //struct of the userData
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
	info, exists := allSessions[string([]byte(sessionKey))]
	if exists {
		delta := info.validUntil.Sub(time.Now())
		if delta > 0 {
			return true
		}
	}
	//session invalid
	return false
}

//function which returns a sessionkey and saves the key on server side in a map, with more details, without base 64
func generateSessionKey(userID int) string {
	sessionKey := getRandomString(sessionKeyLen)
	expiresOn := time.Now().Add(time.Hour * sessionKeyExpires)
	allSessions[sessionKey] = sessionKeyInfo{
		validUntil: expiresOn,
		forUser:    userID,
	}
	return sessionKey
	//return base64.StdEncoding.EncodeToString([]byte(sessionKey))
}

//deletes a sessionkey / makes it invalid
func delSessionKey(sessionKey string) {
	delete(allSessions, sessionKey)
}

//Function which returns a String of a costum length from a charset
func getRandomString(keyLen int) string {
	//rand.Seed(time.Now().UnixNano())
	var charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789" //random charset for random number generator
	var generatedKey string
	for i := 0; i <= keyLen; i++ {
		generatedKey += string(charset[rand.Intn(len(charset))]) //with number it picks random char from charset
	}
	return generatedKey
}

//Function to Login the user and return 1) if it was successfull 2) an error Message
func login(userName string, password string) (bool, string) {
	userDataFile, error := os.OpenFile(dbLocation, os.O_RDONLY|os.O_CREATE, 0666)
	defer userDataFile.Close()
	if error == nil {
		reader := csv.NewReader(userDataFile)
		for {
			line, err := reader.Read()
			if err == nil {
				if line[1] == userName {
					if comparePasswords(password, line[3]) {
						userID, _ := strconv.Atoi(line[0])
						return true, generateSessionKey(userID)
					} else {
						return false, ErrorMessageLoginPasswordWrong
					}
				}
			} else {
				break
			}
		}
		return false, ErrorMessageLoginUsernameUnknown
	}
	return false, ErrorUnknown
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

//the register function which checks if user already exists and calls the function which appends new user to file
func register(userName string, email string, password string, confirmPass string) (bool, string) {
	if password != confirmPass {
		return false, ErrorMessageRegisterNotSamePass
	} else if len(password) < 8 {
		return false, ErrorMessageRegisterInvalidPasswordPolicy
	}
	readDB()
	//delete symbols that break csv
	userName = strings.Replace(userName, ",", "", -1)
	email = strings.Replace(email, ",", "", -1)
	var maxID int = 0
	for k, v := range userDataMap {
		if userName == v.userName {
			return false, ErrorMessageRegisterUsernameTaken
		}
		if k > maxID { //find out the highest userID to set the new one higher than the highest
			maxID = k
		}
	}
	newUser := userData{maxID + 1, userName, email, hashPassword(password)}
	userDataMap[maxID+1] = newUser
	if appendToDB(newUser) {
		return true, generateSessionKey(newUser.uID)
	} else {
		return false, ErrorUnknown
	}

}

func appendToDB(user userData) bool { //appends the new User to the DB/csv file
	var newline string = strconv.Itoa(user.uID) + "," + user.userName + "," + user.email + "," + user.password + "\n"
	userDataFile, err := os.OpenFile(dbLocation, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer userDataFile.Close()
	if err != nil {
		log.Fatal(err)
		return false
	}
	_, err = userDataFile.WriteString(newline)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
func readDB() { //reads the csv and saves in the map and runtime
	userDataMap = make(map[int]userData)
	userDataFile, err := os.OpenFile(dbLocation, os.O_RDONLY|os.O_CREATE, 0666)
	defer userDataFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(userDataFile)
	for scanner.Scan() {
		user := strings.Split(scanner.Text(), ",")
		if len(user) != 4 {
			continue //invalid data or  empty row
		}
		userID, _ := strconv.Atoi(user[0])
		newUser := userData{userID, user[1], user[2], user[3]}
		userDataMap[userID] = newUser
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func dropTable() { //for testing       https://stackoverflow.com/questions/44416645/truncate-a-file-in-golang
	/*err := os.Remove(dbLocation)
	if err != nil {
		fmt.Println("cannot detele file:" + err.Error())
	}
	_, err2 := os.Create(dbLocation)
	if err2 != nil {
		fmt.Println("cannot create file:" + err2.Error())
	}*/
	userDataFile, err := os.OpenFile(dbLocation, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error " + err.Error())
	}
	defer userDataFile.Close()
	userDataFile.Truncate(0)
	userDataFile.Seek(0, 0)
}
