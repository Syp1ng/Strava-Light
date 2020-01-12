package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	dbLocation = "../../DataStorage/UserDataDB.csv"

	dropTable()
}
func beforeTestLoginData() {
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	dropTable()
}

var actualPassword = "TestPassword"
var hashedPass string

func TestHashPassword(t *testing.T) {
	//test if password is not same after hashfunction
	hashedPass = hashPassword(actualPassword)
	assert.NotEqual(t, actualPassword, hashedPass, "Should be hashed")
}
func TestComparePasswords(t *testing.T) {
	//test if password comparison is working
	hashedPass := hashPassword(actualPassword)
	assert.True(t, comparePasswords(actualPassword, hashedPass), "Compare function of passwords should return true")
}
func TestGetRandomString(t *testing.T) {
	//Test Randomness
	random1 := getRandomString(10)
	random2 := getRandomString(10)
	assert.NotEqual(t, random1, random2, "2 Random Keys shouldn't be equal")
}
func TestRegistration(t *testing.T) {
	//Test Registration
	beforeTestLoginData()

	assert.Equal(t, getUID("111111111111111111111111111"), 0, "it should be 0, because nobody with UserID 0")
	worked, errorMessage := register("testUser", "testUser@users.de", "notEqual", "not")
	assert.False(t, worked, "Should be false, because Password not equal")
	assert.Equal(t, errorMessage, ErrorMessageRegisterNotSamePass, "The Error Message should be displayed")
	worked, errorMessage = register("testUser", "testUser@users.de", "pass", "pass")
	assert.False(t, worked, "Should be false, because Password too short")
	assert.Equal(t, errorMessage, ErrorMessageRegisterInvalidPasswordPolicy, "The Error Message should be displayed")
	worked, errorMessage = register("testUser", "testUser@users.de", "password123", "password123")
	assert.True(t, worked, "Should be true, because valid creds")
	worked, errorMessage = register("testUser", "testUser@users.de", "password123", "password123")
	assert.False(t, worked, "Should be false, because already exists")
	assert.Equal(t, errorMessage, ErrorMessageRegisterUsernameTaken, "The Error Message should be displayed")

}
func TestLogin(t *testing.T) { //login test with Session Key and getUID test
	beforeTestLoginData()

	//Test invalid login
	worked, errorMessage := login("testUser", "password123")
	assert.False(t, worked, "Should be false, because user not registered")
	assert.Equal(t, errorMessage, ErrorMessageLoginUsernameUnknown, "The Error Message should be displayed")
	//register
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")
	//Test valid Login*/
	var sessionKey string
	worked, sessionKey = login("testUser", "password123")
	assert.True(t, worked, "Should be true, login success")
	assert.True(t, checkSessionKey(sessionKey), "Key should be valid")
	assert.True(t, getUID(sessionKey) > 0, "When Session Key valid, it should return uID >0")
	assert.True(t, true, true, "true")
}
