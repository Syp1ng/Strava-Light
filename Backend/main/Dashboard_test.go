package main

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	dbLocationActivity = "../../DataStorage/ActivityDB.csv"
}

func TestEditActivity(t *testing.T) {
	dropTable()
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
		UserID:      5,
		Activityart: "TestActivity",
	}
	assert.False(t, activityMap[1].Comment == "geaendert", "should be the new value")
	assert.False(t, activityMap[1].Activityart == "TestActivity", "should be the new value")
	editActivity(editedAct)
	readAcivityDB()
	println(activityMap[1].Comment)
	assert.True(t, activityMap[1].Comment == "geaendert", "should be the new value")
	assert.True(t, activityMap[1].Activityart == "TestActivity", "should be the new value")
}
func TestGetDataForUser(t *testing.T) {

}
func TestRemoveActivity(t *testing.T) {
	dropTable()
	standardAct := Activity{
		ActID:       1,
		Comment:     "comment",
		UserID:      1,
		Activityart: "Laufen",
	}
	appendToDBACT(standardAct)
	assert.True(t, activityMap[1].Comment == "geaendert", "should be the new value")
	assert.True(t, activityMap[1].Activityart == "TestActivity", "should be the new value")
	removeActivity(1, 1)
	readAcivityDB()
	assert.True(t, len(activityMap) == 0, "file should now be empty")
}
func TestSaveNewData(t *testing.T) {

}
func TestSearch(t *testing.T) {

}
func TestDownloadActivity(t *testing.T) {

}
