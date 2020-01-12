//////////5422223//////////9872387//////////8190324
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Activity struct { //Aufbau einer Aktivität
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
	Timestamp      string
	ZipName        string
}

var dbLocationActivity = "DataStorage/ActivityDB.csv" //Zentraler Verweis auf den Speicherort der Aktivitäten in einer CSV Datei
var activityMap map[int]Activity

func uploadfile(zipname string, filename string, activity string, kommentar string, uid int) {
	readAcivityDB() //Funktion liest alle Aktivitäten in der CSV File aus und speichert sie in einer Map
	maxID := 0
	for _, j := range activityMap { //Map wird durchlaufen um höchste ActID herrauszufinden
		if j.ActID > maxID {
			maxID = j.ActID
		}
	}

	//Definition einer neuen Aktivität (zuerst mit Standart Werten, welche alle in der folgenden Funktion überschrieben werden sollen
	newAct := Activity{maxID + 1, uid, filename, activity, kommentar, 0.0, 0.0, 0.0, "", 0.0, 0, 0.0, 0, 1000, "", zipname}

	//Funktion, die die hochgeladene Datei auswertet
	//Die Datei ist über den filename aufrufbar
	newAct = parseDoc(newAct)

	//Überprüfung eines offensichtlichen Eingabefehlers
	//Wenn Geschwindigkeit über 5 M/s ist und Laufen ausgewählt wurde, wird die Art automatisch zu Radfahren gesetzt
	if newAct.Activityart == "Laufen" && newAct.Avgspeed > 5.0 {
		newAct.Activityart = "Radfahren"
	}
	//Wenn Geschwindigkeit unter 3 M/s ist und Radfahren ausgewählt wurde, wird die Art automatisch zu Laufen gesetzt
	if newAct.Activityart == "Radfahren" && newAct.Avgspeed < 3.0 {
		newAct.Activityart = "Laufen"
	}
	//Funktion, die die Aktivität in die CSV Datei schreibt
	appendToDBACT(newAct)

}

func appendToDBACT(act Activity) bool { //Funktion, die die Aktivität in die CSV Datei schreibt
	//Definieren der neuen Zeile in der CSV Datei im String Format
	var newline = strconv.Itoa(act.ActID) + "," + strconv.Itoa(act.UserID) + "," + act.Filename + "," +
		act.Activityart + "," + act.Comment + "," + fmt.Sprintf("%f", act.Distance) + "," + fmt.Sprintf("%f", act.Standzeit) + "," +
		fmt.Sprintf("%f", act.HighSpeed) + "," + act.Highspeedtime + "," + fmt.Sprintf("%f", act.Avgspeed) +
		"," + strconv.Itoa(act.AvgSpeedFastKM) + "," + fmt.Sprintf("%f", act.AvgSpeedFastMS) + "," +
		strconv.Itoa(act.AvgSpeedSlowKM) + "," + fmt.Sprintf("%f", act.AvgSpeedSlowMS) + "," + act.Timestamp + "," + act.ZipName +
		"\n"
	f, err := os.OpenFile(dbLocationActivity, os.O_APPEND|os.O_WRONLY, os.ModeAppend) //Zugriff auf die Datei, um zu checken ob kein Fehler auftritt
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.WriteString(newline) //Definierte Zeile mit Werten in Datei schreiben
	if err != nil {
		return false
	}
	return true
}

func readAcivityDB() { //Funktion liest alle Aktivitäten in der CSV File aus und speichert sie in einer Map
	var id = 0
	file, err := os.Open(dbLocationActivity) //Zugriff auf die Datei, um zu checken ob kein Fehler auftritt
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	activityMap = nil //Map reseten
	activityMap = make(map[int]Activity)
	scanner := bufio.NewScanner(file) //Scanner der Zeile für Zeile durch die CSV Datei läuft
	for scanner.Scan() {
		activity := strings.Split(scanner.Text(), ",") //Trennen einer Zeile in verschiedene Segmente durch Komma
		actID, err := strconv.Atoi(activity[0])        //Schreiben der einzelnen Segmente in die dazugehörige Position
		userID, err := strconv.Atoi(activity[1])       //der Aktivität
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

		//Schreiben der erstellen Aktivität in die Map
		newActivity := Activity{actID, userID, activity[2], activity[3], activity[4], distance, standzeit, highSpeed, activity[8], avgspeed, avgSpeedFastKM, avgSpeedFastMS, avgSpeedSlowKM, avgSpeedSlowMS, activity[14], activity[15]}
		activityMap[id] = newActivity
		id++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
