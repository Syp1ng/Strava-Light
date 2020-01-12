package main

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func init() {
	dbLocationActivity = "../../DataStorage/ActivityDB.csv"
	tempFilePath = "../../DataStorage/Temp.csv"
	backUpPath = "../../DataStorage/BackupActivityDB.csv"
	activityMap = make(map[int]Activity)
	activityMapForUser = make(map[int]Activity)
}

func TestEditActivity(t *testing.T) {
	DropActivityData()
	standardAct := Activity{
		ActID:       1,
		Comment:     "comment",
		UserID:      1,
		Activityart: "Laufen",
	}
	appendToDBACT(standardAct)
	editedAct := Activity{
		ActID:       1,
		Comment:     "geaendert",
		UserID:      1,
		Activityart: "TestActivity",
	}
	assert.True(t, activityMap[1].Comment == "comment", "should be the new value")
	assert.True(t, activityMap[1].Activityart == "Laufen", "should be the new value")
	editActivity(editedAct)
	readAcivityDB()
	assert.True(t, activityMap[1].Comment == "geaendert", "should be the new value")
	assert.True(t, activityMap[1].Activityart == "TestActivity", "should be the new value")
}
func TestGetDataForUser(t *testing.T) {
	DropActivityData()
	firstActivity := Activity{
		ActID:       1,
		Comment:     "von User 1",
		UserID:      1,
		Activityart: "Laufen",
	}
	appendToDBACT(firstActivity)
	secondActivity := Activity{
		ActID:       2,
		Comment:     "von User 2",
		UserID:      2,
		Activityart: "TestActivity",
	}
	appendToDBACT(secondActivity)
	userDataMap := getDataForUser(1)
	assert.Equal(t, "von User 1", userDataMap[0].Comment, "should be the initalized value and should be displayed")
	for _, k := range userDataMap {
		assert.NotEqual(t, "von User 2", k.Comment, "activity from user 2 should not be displayed to user 1")
	}
}

func TestRemoveActivity(t *testing.T) {
	DropActivityData()
	standardAct := Activity{
		ActID:       1,
		Comment:     "comment",
		UserID:      1,
		Activityart: "Laufen",
	}
	appendToDBACT(standardAct)
	readAcivityDB()
	assert.True(t, len(activityMap) > 0, "should be something in the activity map")
	assert.True(t, activityMap[1].Comment == "comment", "should be the new value")
	assert.True(t, activityMap[1].Activityart == "Laufen", "should be the new value")
	removeActivity(1, 1)
	readAcivityDB()
	assert.True(t, len(activityMap) == 0, "file should now be empty")
}

func TestSearch(t *testing.T) { //for testing
	DropActivityData()
	firstAct := Activity{
		ActID:       1,
		Comment:     "test",
		UserID:      1,
		Activityart: "Laufen",
	}
	secondAct := Activity{
		ActID:       1,
		Comment:     "don't find me in search",
		UserID:      5,
		Activityart: "Fahrrad",
	}
	appendToDBACT(firstAct)
	appendToDBACT(secondAct)
	testSearch := "test"
	x := search(1, testSearch)
	for _, v := range x {
		assert.True(t, strings.Contains(v.Comment, testSearch))
	}
}
func TestDownloadActivity(t *testing.T) {

}
