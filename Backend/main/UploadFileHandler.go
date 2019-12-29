package main

import (
	"fmt"
	"os"
	"strconv"
)

type Activity struct {
	ActID          int
	UserID         string
	Filename       string
	Activityart    string
	Comment        string
	Distance       float64
	Standzeit      float64
	HighSpeed      float64
	Highspeedtime  string
	Avgspeed       float64
	AvgSpeedFastKM int
	AvgSpeedFastMS float64
	AvgSpeedSlowKM int
	AvgSpeedSlowMS float64
}

var dbLocationActivity = "DataStorage/ActivityDB.csv"

func uploadfile(filename string, activity string, kommentar string, uid string) {

	newAct := Activity{1, uid, filename, activity, kommentar, 0.0, 0.0, 0.0, "", 0.0, 0, 0.0, 0, 1000}
	newAct = parseDoc(newAct)
	if newAct.Activityart == "Laufen" && newAct.Avgspeed > 5.0 {
		newAct.Activityart = "Radfahren"
	}
	if newAct.Activityart == "Radfahren" && newAct.Avgspeed < 3.0 {
		newAct.Activityart = "Laufen"
	}
	appendToDBACT(newAct)

}

func appendToDBACT(act Activity) bool {

	var newline = strconv.Itoa(act.ActID) + "," + act.UserID + "," + act.Filename + "," +
		act.Activityart + "," + act.Comment + "," + fmt.Sprintf("%f", act.Distance) + "," + fmt.Sprintf("%f", act.Standzeit) + "," +
		fmt.Sprintf("%f", act.HighSpeed) + "," + act.Highspeedtime + "," + fmt.Sprintf("%f", act.Avgspeed) +
		"," + strconv.Itoa(act.AvgSpeedFastKM) + "," + fmt.Sprintf("%f", act.AvgSpeedFastMS) + "," +
		strconv.Itoa(act.AvgSpeedSlowKM) + "," + fmt.Sprintf("%f", act.AvgSpeedSlowMS) +
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
