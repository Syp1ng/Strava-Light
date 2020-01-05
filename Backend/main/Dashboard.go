package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var activityMapForUser map[int]Activity

func getDataForUser(uID int) map[int]Activity {
	var position = 0
	readAcivityDB()
	activityMapForUser = nil
	activityMapForUser = make(map[int]Activity)
	for _, j := range activityMap {
		if j.UserID == uID {
			userActivity := j
			activityMapForUser[position] = userActivity
			position++
		}
	}
	return activityMapForUser
}

func removeActivity(uID int, activityID int) {
	file, err := os.Open(dbLocationActivity)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	f, err := os.OpenFile("DataStorage/Temp.csv", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {

	}
	defer f.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		activity := strings.Split(scanner.Text(), ",")
		actID, err := strconv.Atoi(activity[0])
		//userID, err := strconv.Atoi(activity[1])
		if err != nil {

		}
		if actID != activityID {
			f.WriteString(scanner.Text() + "\n")
		} else {
			var error = os.Remove(activity[2])
			if error != nil {
				log.Println("Can't find the file ", error)
			}
		}
	}
	file.Close()
	f.Close()
	var error = os.Rename(filepath.Join("DataStorage", "ActivityDB.csv"), filepath.Join("DataStorage", "BackupActivityDB.csv"))
	if error != nil {
		log.Println(error)
	} else {
		error = os.Rename(filepath.Join("DataStorage", "Temp.csv"), filepath.Join("DataStorage", "ActivityDB.csv"))
		if error != nil {
			log.Println(error)
		}
	}
	//os.Remove(path)
}

func editActivity(editactivity Activity) {
	file, err := os.Open(dbLocationActivity)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	f, err := os.OpenFile("DataStorage/Temp.csv", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {

	}
	defer f.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		actID, err := strconv.Atoi(line[0])
		//userID, err := strconv.Atoi(activity[1])
		if err != nil {

		}
		if actID != editactivity.ActID {
			f.WriteString(scanner.Text() + "\n")

		} else {
			editactivity.Filename = line[2]
			editactivity.Distance, err = strconv.ParseFloat(line[5], 64)
			editactivity.Standzeit, err = strconv.ParseFloat(line[6], 64)
			editactivity.HighSpeed, err = strconv.ParseFloat(line[7], 64)
			editactivity.Highspeedtime = line[8]
			editactivity.Avgspeed, err = strconv.ParseFloat(line[9], 64)
			editactivity.AvgSpeedFastKM, err = strconv.Atoi(line[10])
			editactivity.AvgSpeedFastMS, err = strconv.ParseFloat(line[11], 64)
			editactivity.AvgSpeedSlowKM, err = strconv.Atoi(line[12])
			editactivity.AvgSpeedSlowMS, err = strconv.ParseFloat(line[13], 64)
		}
	}
	var newline = strconv.Itoa(editactivity.ActID) + "," + strconv.Itoa(editactivity.UserID) + "," + editactivity.Filename + "," +
		editactivity.Activityart + "," + editactivity.Comment + "," + fmt.Sprintf("%f", editactivity.Distance) + "," + fmt.Sprintf("%f", editactivity.Standzeit) + "," +
		fmt.Sprintf("%f", editactivity.HighSpeed) + "," + editactivity.Highspeedtime + "," + fmt.Sprintf("%f", editactivity.Avgspeed) +
		"," + strconv.Itoa(editactivity.AvgSpeedFastKM) + "," + fmt.Sprintf("%f", editactivity.AvgSpeedFastMS) + "," +
		strconv.Itoa(editactivity.AvgSpeedSlowKM) + "," + fmt.Sprintf("%f", editactivity.AvgSpeedSlowMS) +
		"\n"
	_, err = f.WriteString(newline)
	if err != nil {
		log.Println(err)
	}

	file.Close()
	f.Close()
	var error = os.Rename(filepath.Join("DataStorage", "ActivityDB.csv"), filepath.Join("DataStorage", "BackupActivityDB.csv"))
	if error != nil {
		log.Println(error)
	} else {
		error = os.Rename(filepath.Join("DataStorage", "Temp.csv"), filepath.Join("DataStorage", "ActivityDB.csv"))
		if error != nil {
			log.Println(error)
		}
	}
}
func saveNewData(activity []Activity) {

}
func search(uID int, comment string) []Activity {
	var commentaryMap map[int]Activity
	commentaryMap = make(map[int]Activity)
	readAcivityDB()
	for key, value := range activityMap {
		if value.UserID == uID && strings.Contains(value.Comment, comment) {
			fmt.Println("Key:", key, "Value:", value)
			commentaryMap[key] = value
		}
	}
	//how to convert map into Acitivtiy Array ? maybe:https://stackoverflow.com/questions/45570947/creating-an-array-from-the-maps-key-and-values-in-go/45571006
	//What to return ?
	return nil
}
