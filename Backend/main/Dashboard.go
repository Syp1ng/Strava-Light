package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func getDataForUser(uID int) []Activity {
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
		return slice
	}
	return nil

	//from dbLocationActivity
	/*var lala = []Activity{
		Activity{
			ActID:          1,
			UserID:         "1",
			filename:       "asdsa",
			activityart:    "sada",
			comment:        "asdasd",
			distance:       5.0,
			standzeit:      5.0,
			highSpeed:      5.0,
			highspeedtime:  "string",
			avgspeed:       5.0,
			avgSpeedFastKM: 3,
			avgSpeedFastMS: 4,
			avgSpeedSlowKM: 5,
			avgSpeedSlowMS: 5.6,
		},
		Activity{
			ActID:          3,
			UserID:         "4",
			filename:       "asdsa",
			activityart:    "sada",
			comment:        "asdasd",
			distance:       10.0,
			standzeit:      5.0,
			highSpeed:      5.0,
			highspeedtime:  "string",
			avgspeed:       5.0,
			avgSpeedFastKM: 3,
			avgSpeedFastMS: 4,
			avgSpeedSlowKM: 5,
			avgSpeedSlowMS: 5.6,
		},
	}*/

}

func removeActivity(uID int, activityID int) {
}
func editActivity(activity Activity) {
	fmt.Println(activity.Comment + " " + activity.Activityart)
}
func saveNewData(activity []Activity) {

}
