package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	dbLocation = "../../DataStorage/UserDataDB.csv"

	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	dropTable()
}
func beforeTest() {
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
	dropTable()
}

var actualPassword = "TestPassword"
var hashedPass string

func TestHashPassword(t *testing.T) {
	hashedPass = hashPassword(actualPassword)
	assert.NotEqual(t, actualPassword, hashedPass, "Should be hashed")
}
func TestComparePasswords(t *testing.T) {
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
	beforeTest()

	assert.Equal(t, getUID("111111111111111111111111111"), 0, "it should be 0, because nobody with UserID 0")
	worked, _ := register("testUser", "testUser@users.de", "notEqual", "not")
	assert.False(t, worked, "Should be false, because Password not equal")
	worked, _ = register("testUser", "testUser@users.de", "pass", "pass")
	assert.False(t, worked, "Should be false, because Password too short")
	worked, _ = register("testUser", "testUser@users.de", "password123", "password123")
	assert.True(t, worked, "Should be true, because valid creds")
	worked, _ = register("testUser", "testUser@users.de", "password123", "password123")
	assert.False(t, worked, "Should be false, because already exists")

}
func TestLogin(t *testing.T) {
	beforeTest()

	//Test invalid login
	worked, _ := login("testUser", "password123")
	assert.False(t, worked, "Should be false, because user not registered")
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
