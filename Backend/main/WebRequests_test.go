package main

import (
	_ "github.com/stretchr/testify"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	_ "testing"
)

//https://stackoverflow.com/questions/30  	978105/adding-post-variables-to-go-test-http-request

func init() {
	registerTemplatePath = "../../Frontend/RegisterTemplate.html"
	loginTemplatePath = "../../Frontend/LoginTemplate.html"
	dbLocation = "../../DataStorage/UserDataDB.csv"
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
}

func TestLogoutHandler(t *testing.T) {
	beforeTest()
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")
	_, sessionKey := login("testUser", "password123")
	assert.True(t, checkSessionKey(sessionKey), "session Key has to be vaild first")

	//Logout
	req, _ := http.NewRequest("POST", "/logout", nil)
	cookie := http.Cookie{Name: "auth", Value: sessionKey}
	req.AddCookie(&cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	logoutHandler(w, req)
	assert.False(t, checkSessionKey(sessionKey), "session Key should be invalid after logout")
}

func TestEditHandler(t *testing.T) {

}
func TestRemoveHandler(t *testing.T) {

}

func TestDownloadHandler(t *testing.T) {

}

func TestViewDashboardHandler(t *testing.T) {

}

func TestRegisterHandler(t *testing.T) {
	beforeTest()

	reader := strings.NewReader("username=testUser&password=password123&confirmPassword=password12&email=testUser@users.de")
	req, _ := http.NewRequest("POST", "/registrationHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	registerHandler(w, req)
	html := w.Body.String()
	assert.True(t, strings.Contains(html, ErrorMessageRegisterNotSamePass), "HTML should contain error message that password is not same")

	reader = strings.NewReader("username=testUser&password=pass&confirmPassword=pass&email=testUser@users.de")
	req, _ = http.NewRequest("POST", "/registrationHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	registerHandler(w, req)
	html = w.Body.String()
	assert.True(t, strings.Contains(html, ErrorMessageRegisterInvalidPasswordPolicy), "HTML should contain error message that user password too short")

	//check valid registration
	reader = strings.NewReader("username=testUser&password=password123&confirmPassword=password123&email=testUser@users.de")
	req, _ = http.NewRequest("POST", "/registrationHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	registerHandler(w, req)
	assert.Equal(t, http.StatusFound, w.Result().StatusCode, "when valid login status is fine and redirect")

	//check invalid registration
	reader = strings.NewReader("username=testUser&password=password123&confirmPassword=password123&email=testUser@users.de")
	req, _ = http.NewRequest("POST", "/registrationHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	registerHandler(w, req)
	html = w.Body.String()
	assert.True(t, strings.Contains(html, ErrorMessageRegisterUsernameTaken), "HTML should contain error message that user is taken")
}

func TestLoginHandler(t *testing.T) {
	beforeTest()
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")

	//Username Unknown
	reader := strings.NewReader("username=testUser1&password=password123")
	req, _ := http.NewRequest("POST", "/loginHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	loginHandler(w, req)
	html := w.Body.String()
	assert.True(t, strings.Contains(html, ErrorMessageLoginUsernameUnknown), "HTML should contain error message that user not exists")

	//Wrong Password
	reader = strings.NewReader("username=testUser&password=password1234")
	req, _ = http.NewRequest("POST", "/loginHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	loginHandler(w, req)
	html = w.Body.String()
	assert.True(t, strings.Contains(html, ErrorMessageLoginPasswordWrong), "HTML should contain error message that password is wrong")

	//valid login
	reader = strings.NewReader("username=testUser&password=password123")
	req, _ = http.NewRequest("POST", "/loginHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	loginHandler(w, req)
	assert.Equal(t, http.StatusFound, w.Result().StatusCode, "when valid login status is fine and redirect")

}

func TestUploadHandler(t *testing.T) {

}
func TestUnzip(t *testing.T) {
}
