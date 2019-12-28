package main

import (
	"fmt"
	"os"
	"strconv"
)

type Activity struct {
	actID          int
	UserID         int
	filename       string
	activityart    string
	comment        string
	distance       float64
	standzeit      float64
	highSpeed      float64
	highspeedtime  string
	avgspeed       float64
	avgSpeedFastKM int
	avgSpeedFastMS float64
	avgSpeedSlowKM int
	avgSpeedSlowMS float64
}

var dbLocationActivity = "DataStorage/ActivityDB.csv"

func uploadfile(filename string, activity string, kommentar string) {

	newAct := Activity{1, 1, filename, activity, kommentar, 0.0, 0.0, 0.0, "", 0.0, 0, 0.0, 0, 1000}
	newAct = parseDoc(newAct)
	if newAct.activityart == "Laufen" && newAct.avgspeed > 5.0 {
		newAct.activityart = "Radfahren"
	}
	if newAct.activityart == "Radfahren" && newAct.avgspeed < 3.0 {
		newAct.activityart = "Laufen"
	}
	appendToDBACT(newAct)

}

func appendToDBACT(act Activity) bool {

	var newline = strconv.Itoa(act.actID) + "," + strconv.Itoa(act.UserID) + "," + act.filename + "," +
		act.activityart + "," + act.comment + "," + fmt.Sprintf("%f", act.distance) + "," + fmt.Sprintf("%f", act.standzeit) + "," +
		fmt.Sprintf("%f", act.highSpeed) + "," + act.highspeedtime + "," + fmt.Sprintf("%f", act.avgspeed) +
		"," + strconv.Itoa(act.avgSpeedFastKM) + "," + fmt.Sprintf("%f", act.avgSpeedFastMS) + "," +
		strconv.Itoa(act.avgSpeedSlowKM) + "," + fmt.Sprintf("%f", act.avgSpeedSlowMS) +
		"\n"
	f, err := os.OpenFile(dbLocationActivity, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
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
