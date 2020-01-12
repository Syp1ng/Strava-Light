//////////5422223//////////9872387//////////8190324//////////
package main

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func init() {
	dbLocationActivity = "../../DataStorage/ActivityDB.csv"
	backUpPath = "../../DataStorage/BackupActivityDB.csv"
	tempFilePath = "../../DataStorage/Temp.csv"
	activityMap = make(map[int]Activity)
	activityMapForUser = make(map[int]Activity)
}
func beforeTestActivityData() {
	DropActivityData()
}
func addExampleActivityForTest() {
	standardAct := Activity{
		Comment:     "comment",
		UserID:      1,
		Activityart: "Laufen",
	}
	appendToDBACT(standardAct)
	readAcivityDB()
}
func TestEditActivity(t *testing.T) {
	beforeTestActivityData()
	addExampleActivityForTest()
	editedAct := Activity{
		ActID:       0,
		Comment:     "geaendert",
		UserID:      1,
		Activityart: "TestActivity",
	}
	assert.True(t, activityMap[0].Comment == "comment", "should be the new value")
	assert.True(t, activityMap[0].Activityart == "Laufen", "should be the new value")
	editActivity(editedAct)
	readAcivityDB()
	assert.True(t, activityMap[0].Comment == "geaendert", "should be the new value")
	assert.True(t, activityMap[0].Activityart == "TestActivity", "should be the new value")
}
func TestGetDataForUser(t *testing.T) {
	DropActivityData()
	firstActivity := Activity{
		ActID:       0,
		Comment:     "von User 1",
		UserID:      1,
		Activityart: "Laufen",
	}
	appendToDBACT(firstActivity)
	secondActivity := Activity{
		ActID:       1,
		Comment:     "von User 2",
		UserID:      2,
		Activityart: "TestActivity",
	}
	appendToDBACT(secondActivity)
	userDataMap := getDataForUser(1)
	assert.True(t, len(userDataMap) == 1, " should contain 1 element, the other activity is for user 2")
	assert.Equal(t, "von User 1", userDataMap[0].Comment, "should be the initalized value and should be displayed")
	for _, k := range userDataMap {
		assert.NotEqual(t, "von User 2", k.Comment, "activity from user 2 should not be displayed to user 1")
	}
}

func TestRemoveActivity(t *testing.T) {
	DropActivityData()
	addExampleActivityForTest()
	assert.True(t, len(activityMap) > 0, "should be something in the activity map")
	for _, v := range activityMap { //only iterate one time when only one Activit is there
		assert.True(t, v.Comment == "comment", "should be the new value")
		assert.True(t, v.Activityart == "Laufen", "should be the new value")
	}
	removeActivity(1, 0)
	readAcivityDB()
	assert.True(t, len(activityMap) == 0, "file should now be empty")
}

func TestSearch(t *testing.T) { //for testing
	DropActivityData()
	addExampleActivityForTest()
	secondAct := Activity{
		ActID:       1,
		Comment:     "don't find me in search",
		UserID:      1,
		Activityart: "Fahrrad",
	}
	appendToDBACT(secondAct)
	testSearch := "test"
	x := search(1, testSearch)
	for _, v := range x {
		assert.True(t, strings.Contains(v.Comment, testSearch))
	}
}
