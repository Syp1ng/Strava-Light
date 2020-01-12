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
	dashboardTemplatePath = "../../Frontend/dashboardTemplate.html"
	dbLocationActivity = "../../DataStorage/ActivityDB.csv"
	dbLocation = "../../DataStorage/UserDataDB.csv"
	allSessions = make(map[string]sessionKeyInfo)
	userDataMap = make(map[int]userData)
}

func TestLogoutHandler(t *testing.T) {
	beforeTest()
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")
	_, sessionKey := login("testUser", "password123")
	assert.True(t, checkSessionKey(sessionKey), "session Key has to be vaild first")

	// test Logout
	req, _ := http.NewRequest("POST", "/logout", nil)
	cookie := http.Cookie{Name: "auth", Value: sessionKey}
	req.AddCookie(&cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	logoutHandler(w, req)
	assert.False(t, checkSessionKey(sessionKey), "session Key should be invalid after logout")
}

func TestEditHandler(t *testing.T) {
	beforeTest()
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")
	_, sessionKey := login("testUser", "password123")

	//invalid creds && valid parameters
	reader := strings.NewReader("actID=1&actArt=Testart&comment=geaendert")
	req, _ := http.NewRequest("POST", "/editActivity", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cookie := http.Cookie{Name: "auth", Value: "11111"}
	req.AddCookie(&cookie)
	w := httptest.NewRecorder()
	editHandler(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, " should not have permissions")

	//valid creds && valid parameters
	reader = strings.NewReader("actID=1&actArt=Testart&comment=geaendert")
	req, _ = http.NewRequest("POST", "/editActivity", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cookie = http.Cookie{Name: "auth", Value: sessionKey}
	req.AddCookie(&cookie)
	w = httptest.NewRecorder()
	editHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode, " should have permissions")
}
func TestRemoveHandler(t *testing.T) {
	beforeTest()
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")
	_, sessionKey := login("testUser", "password123")

	//invalid creds && valid parameters
	reader := strings.NewReader("actID=1&actArt=Testart&comment=geaendert")
	req, _ := http.NewRequest("POST", "/removeActivity", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cookie := http.Cookie{Name: "auth", Value: "11111"}
	req.AddCookie(&cookie)
	w := httptest.NewRecorder()
	removeHandler(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, " should not have permissions")

	//valid creds && valid parameters
	reader = strings.NewReader("actID=1&actArt=Testart&comment=geaendert")
	req, _ = http.NewRequest("POST", "/removeActivity", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cookie = http.Cookie{Name: "auth", Value: sessionKey}
	req.AddCookie(&cookie)
	w = httptest.NewRecorder()
	removeHandler(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode, " should have permissions")
}

func TestDownloadHandler(t *testing.T) {

}

func TestViewDashboardHandler(t *testing.T) {
	beforeTest()
	//with no access
	req, _ := http.NewRequest("POST", "/home", nil)
	cookie := http.Cookie{Name: "auth", Value: "11111111"}
	req.AddCookie(&cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	viewDashboardHandler(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode, "should not have to access to dashboard site. Error with unauthorized")

	//with valid access
	_, _ = register("testUser", "testUser@users.de", "password123", "password123")
	_, sessionKey := login("testUser", "password123")

	// test if you get dashboard
	req, _ = http.NewRequest("POST", "/home", nil)
	cookie = http.Cookie{Name: "auth", Value: sessionKey}
	req.AddCookie(&cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	viewDashboardHandler(w, req)
	html := w.Body.String()
	assert.True(t, strings.Contains(html, "<title>Dashboard</title>") && strings.Contains(html, "<h1>Neue Aktivität hochladen</h1>"), "Dashboard should bedisplayed")
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
	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "status has to be ok")

	//Wrong Password
	reader = strings.NewReader("username=testUser&password=password1234")
	req, _ = http.NewRequest("POST", "/loginHandler", reader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()
	loginHandler(w, req)
	html = w.Body.String()
	assert.True(t, strings.Contains(html, ErrorMessageLoginPasswordWrong), "HTML should contain error message that password is wrong")
	assert.Equal(t, http.StatusOK, w.Result().StatusCode, "status has to be ok")

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
