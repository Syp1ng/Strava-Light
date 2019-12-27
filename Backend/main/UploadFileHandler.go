package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Activity struct {
	actID         int
	UserID        int
	filename      string
	activityart   string
	comment       string
	distance      float64
	highSpeed     float64
	highspeedtime string
	avgspeed      float64
}

var dbLocationActivity = "DataStorage/ActivityDB.csv"

func uploadfile(filename string, activity string, kommentar string) {

	log.Println(filename)
	newAct := Activity{1, 1, filename, activity, kommentar, 0.0, 0.0, "", 0.0}
	log.Println(newAct.distance)
	newAct = parseDoc(newAct)
	log.Println(newAct.distance)
	appendToDBACT(newAct)

}

func appendToDBACT(act Activity) bool {

	var newline string = strconv.Itoa(act.actID) + "," + strconv.Itoa(act.UserID) + "," + act.filename + "," + act.activityart + "," + act.comment + "," + fmt.Sprintf("%f", act.distance) + "," + fmt.Sprintf("%f", act.highSpeed) + "," + act.highspeedtime + "," + fmt.Sprintf("%f", act.avgspeed) + "\n"
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
