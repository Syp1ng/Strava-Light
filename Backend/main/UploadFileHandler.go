package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Activity struct {
	ActID          int
	UserID         int
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
var activityMap map[int]Activity

func uploadfile(filename string, activity string, kommentar string, uid int) {
	readAcivityDB()
	maxID := 0
	for k := range activityMap {
		if k > maxID {
			maxID = k
		}

	}

	newAct := Activity{maxID + 1, uid, filename, activity, kommentar, 0.0, 0.0, 0.0, "", 0.0, 0, 0.0, 0, 1000}
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

	var newline = strconv.Itoa(act.ActID) + "," + strconv.Itoa(act.UserID) + "," + act.Filename + "," +
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

func readAcivityDB() {
	file, err := os.Open(dbLocationActivity)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		activity := strings.Split(scanner.Text(), ",")
		actID, err := strconv.Atoi(activity[0])
		userID, err := strconv.Atoi(activity[1])
		distance, err := strconv.ParseFloat(activity[5], 64)
		standzeit, err := strconv.ParseFloat(activity[6], 64)
		highSpeed, err := strconv.ParseFloat(activity[7], 64)
		avgspeed, err := strconv.ParseFloat(activity[9], 64)
		avgSpeedFastKM, err := strconv.Atoi(activity[10])
		avgSpeedFastMS, err := strconv.ParseFloat(activity[11], 64)
		avgSpeedSlowKM, err := strconv.Atoi(activity[12])
		avgSpeedSlowMS, err := strconv.ParseFloat(activity[13], 64)
		if err != nil {
			log.Fatal(err)
		}

		newActivity := Activity{actID, userID, activity[2], activity[3], activity[3], distance, standzeit, highSpeed, activity[7], avgspeed, avgSpeedFastKM, avgSpeedFastMS, avgSpeedSlowKM, avgSpeedSlowMS}
		activityMap[actID] = newActivity
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func searchInComment(searching string) {
	readAcivityDB()
	fmt.Println(strings.Contains("Assume", "Ass"))
}
