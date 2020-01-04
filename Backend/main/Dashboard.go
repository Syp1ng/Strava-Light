package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getDataForUser(uID int) []Activity {
	/*readAcivityDB()
	var length = 0
	for length = range activityMap {
		length++
	}


	*/
	numberoflines := 0
	activityDataLenght, error := os.Open(dbLocationActivity)
	if error == nil {
		readerforlengh := csv.NewReader(activityDataLenght)
		for {
			lin, err := readerforlengh.Read()
			if err == nil {
				if lin[1] == strconv.Itoa(uID) {
					numberoflines++
				}

			} else {
				break
				log.Println(lin)
			}
		}
	}
	activityDataLenght.Close()
	activityData, err := os.Open(dbLocationActivity)
	if err == nil {
		reader := csv.NewReader(activityData)
		var slice = make([]Activity, numberoflines)
		lines := 0
		for {
			line, err := reader.Read()
			if err == nil {
				if line[1] == strconv.Itoa(uID) {
					slice[lines].ActID, err = strconv.Atoi(line[0])
					slice[lines].UserID, err = strconv.Atoi(line[1])
					slice[lines].Filename = line[2]
					slice[lines].Activityart = line[3]
					slice[lines].Comment = line[4]
					slice[lines].Distance, err = strconv.ParseFloat(line[5], 64)
					slice[lines].Standzeit, err = strconv.ParseFloat(line[6], 64)
					slice[lines].HighSpeed, err = strconv.ParseFloat(line[7], 64)
					slice[lines].Highspeedtime = line[8]
					slice[lines].Avgspeed, err = strconv.ParseFloat(line[9], 64)
					slice[lines].AvgSpeedFastKM, err = strconv.Atoi(line[10])
					slice[lines].AvgSpeedFastMS, err = strconv.ParseFloat(line[11], 64)
					slice[lines].AvgSpeedSlowKM, err = strconv.Atoi(line[12])
					slice[lines].AvgSpeedSlowMS, err = strconv.ParseFloat(line[13], 64)

					lines++
				}

			} else {
				break
			}
		}
		activityData.Close()
		return slice
	}
	activityData.Close()
	return nil

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
			//log.Println("Not allowed to delete")
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
	return nil
}
func downloadActivity(actID int) {

}
