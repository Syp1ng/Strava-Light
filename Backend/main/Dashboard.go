package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var activityMapForUser map[int]Activity
var SortetactivityMapForUser map[int]Activity
var commentaryMap map[int]Activity
var tempFilePath = "DataStorage/Temp.csv"
var backUpPath = "DataStorage/BackupActivityDB.csv"

func DropActivityData() { //like in RegisterAndSessionHandler  	for testinghttps://stackoverflow.com/questions/44416645/truncate-a-file-in-golang
	activityDB, err := os.OpenFile(dbLocationActivity, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error " + err.Error())
	}
	defer activityDB.Close()
	activityDB.Truncate(0)
	activityDB.Seek(0, 0)

	backUpDP, err := os.OpenFile(backUpPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("error " + err.Error())
	}
	defer backUpDP.Close()
	backUpDP.Truncate(0)
	backUpDP.Seek(0, 0)
}

func getDataForUser(uID int) map[int]Activity {
	var position = 0
	readAcivityDB()

	//Map mit Aktivitäten eines Nutzeres reseten
	activityMapForUser = nil
	activityMapForUser = make(map[int]Activity)

	//Nur die Aktivitätetn des Nutzers reinladen, der den Button bedrückt hat
	for _, j := range activityMap {
		if j.UserID == uID {
			userActivity := j
			activityMapForUser[position] = userActivity
			position++

		}
	}
	/*keys := make([]string, 0, len(activityMapForUser))

	for _,k := range activityMapForUser {
		keys = append(keys, k.Timestamp)
	}
	sort.Strings(keys)
	SortetactivityMapForUser = nil
	SortetactivityMapForUser = make(map[int]Activity)

	for _,k := range keys {
		fmt.Println(k)
		//SortetactivityMapForUser[i] = activityMapForUser[k]
	}
	*/
	return activityMapForUser
}

func removeActivity(uID int, activityID int) {
	file, err := os.Open(dbLocationActivity)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//Temporäre Datei erstellen
	f, err := os.OpenFile(tempFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {

	}
	defer f.Close()
	scanner := bufio.NewScanner(file) //Scanner der Zeile für Zeile durch die CSV Datei läuft
	for scanner.Scan() {
		activity := strings.Split(scanner.Text(), ",") //Trennen einer Zeile in verschiedene Segmente durch Komma
		actID, err := strconv.Atoi(activity[0])        //AktivitätsId der aktuellen Aktivität auslesen
		//userID, err := strconv.Atoi(activity[1])
		if err != nil {

		}
		if actID != activityID { //Überprüfen ob die aktuelle Zeile die gesuchte ID enthält
			f.WriteString(scanner.Text() + "\n") //Falls nein, Zeile in Tempräre Datei schreiben
		} else { //Falls ja, Zeile nicht in Temporäre Datei schreiben
			var error = os.Remove(activity[2]) //Verlinkte Datei aus dem DataStorage entfernen
			if error != nil {
				log.Println("Can't find the file ", error)
			}
		}
	}
	file.Close()
	f.Close()

	//Temporäre Datei wird zur Hauptaktivitäts-Datei und die vorherige Hauptaktivitäts-Datei wird zum Backup durch Rename
	var error = os.Rename(dbLocationActivity, backUpPath)
	if error != nil {
		log.Println(error)
	} else {
		error = os.Rename(tempFilePath, dbLocationActivity)
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
	//Temporäre Datei erstellen
	f, err := os.OpenFile(tempFilePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {

	}
	defer f.Close()
	scanner := bufio.NewScanner(file) //Scanner der Zeile für Zeile durch die CSV Datei läuft
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",") //Trennen einer Zeile in verschiedene Segmente durch Komma
		actID, err := strconv.Atoi(line[0])        //AktivitätsId der aktuellen Aktivität auslesen
		//userID, err := strconv.Atoi(activity[1])
		if err != nil {

		}
		if actID != editactivity.ActID { //Überprüfen ob die aktuelle Zeile die gesuchte ID enthält
			f.WriteString(scanner.Text() + "\n") //Falls nein, Zeile in Tempräre Datei schreiben

		} else { //Falls ja, neue Zeile mit bearbeiteten Werten am Ende der CSV Datei anhängen

			editactivity.Filename = line[2] //Auslesen der nicht veränderbaren Werte
			editactivity.Distance, err = strconv.ParseFloat(line[5], 64)
			editactivity.Standzeit, err = strconv.ParseFloat(line[6], 64)
			editactivity.HighSpeed, err = strconv.ParseFloat(line[7], 64)
			editactivity.Highspeedtime = line[8]
			editactivity.Avgspeed, err = strconv.ParseFloat(line[9], 64)
			editactivity.AvgSpeedFastKM, err = strconv.Atoi(line[10])
			editactivity.AvgSpeedFastMS, err = strconv.ParseFloat(line[11], 64)
			editactivity.AvgSpeedSlowKM, err = strconv.Atoi(line[12])
			editactivity.AvgSpeedSlowMS, err = strconv.ParseFloat(line[13], 64)
			editactivity.Timestamp = line[14]
			editactivity.ZipName = line[15]
		}
	}

	//Schreiben aller Werte in neue Zeile der CSV Datei
	var newline = strconv.Itoa(editactivity.ActID) + "," + strconv.Itoa(editactivity.UserID) + "," + editactivity.Filename + "," +
		editactivity.Activityart + "," + editactivity.Comment + "," + fmt.Sprintf("%f", editactivity.Distance) + "," + fmt.Sprintf("%f", editactivity.Standzeit) + "," +
		fmt.Sprintf("%f", editactivity.HighSpeed) + "," + editactivity.Highspeedtime + "," + fmt.Sprintf("%f", editactivity.Avgspeed) +
		"," + strconv.Itoa(editactivity.AvgSpeedFastKM) + "," + fmt.Sprintf("%f", editactivity.AvgSpeedFastMS) + "," +
		strconv.Itoa(editactivity.AvgSpeedSlowKM) + "," + fmt.Sprintf("%f", editactivity.AvgSpeedSlowMS) + "," + editactivity.Timestamp + "," + editactivity.ZipName +
		"\n"
	_, err = f.WriteString(newline)
	if err != nil {
		log.Println(err)
	}

	file.Close()
	f.Close()

	//Temporäre Datei wird zur Hauptaktivitäts-Datei und die vorherige Hauptaktivitäts-Datei wird zum Backup durch Rename
	var error = os.Rename(dbLocationActivity, backUpPath)
	if error != nil {
		log.Println(error)
	} else {
		error = os.Rename(tempFilePath, dbLocationActivity)
		if error != nil {
			log.Println(error)
		}
	}
}

//FullTextSearch
func search(uID int, comment string) map[int]Activity {
	//Zaehler für Index der neuen Map
	var position = 0
	commentaryMap = make(map[int]Activity)
	//ActivityDB auslesen
	readAcivityDB()
	//für alle Aktivitäten in der Map überprüfe
	for _, value := range activityMap {
		//wenn UserId übereinstimmt und das gespeicherte Kommentar den Suchstring beinhaltet dann schreibe in neue Map
		if value.UserID == uID && strings.Contains(value.Comment, comment) {
			commentaryMap[position] = value
			position++
		}
	}
	//gebe neue Map zurück
	return commentaryMap
}
