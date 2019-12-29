package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessions(t *testing.T) {
	//Test Randomness
	assert.NotEqual(t, getRandomString(10), getRandomString(10), "2 Random Keys shouldn't be equal")

	//Test Passwords
	actualPassword := "TestPassword"
	hashedPass := hashPassword(actualPassword)
	assert.NotEqual(t, actualPassword, hashedPass, "Should be hashed")
	assert.True(t, comparePasswords(actualPassword, hashedPass), "Compare function of passwords should return true")

	//Test invalid login
	worked, _ := login("testUser", "password123")
	assert.False(t, worked, "Should be false, because user not registered")

	//Test Registration
	assert.Equal(t, getUID("111111111111111111111111111"), 0, "it should be 0, because nobody with UserID 0")
	worked, _ = register("testUser", "testUser@users.de", "notEqual", "not")
	assert.False(t, worked, "Should be false, because Password not equal")
	worked, _ = register("testUser", "testUser@users.de", "pass", "pass")
	assert.False(t, worked, "Should be false, because Password too short")
	worked, _ = register("testUser", "testUser@users.de", "password123", "password123")
	assert.True(t, worked, "Should be true, because valid creds")
	worked, _ = register("testUser", "testUser@users.de", "password123", "password123")
	assert.False(t, worked, "Should be false, because already exists")

	//Test Login
	var sessionKey string
	worked, sessionKey = login("testUser", "password123")
	assert.False(t, worked, "Should be true, login success")
	assert.True(t, checkSessionKey(sessionKey), "Key should be valid")
	assert.True(t, getUID(sessionKey) > 0, "When Session Key valid, it should return uID >0")
}
