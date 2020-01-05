package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func SetupLinks() {
	http.HandleFunc("/home", viewDashboardHandler)
	http.HandleFunc("/registrationHandler", registerHandler)
	http.HandleFunc("/loginHandler", loginHandler)
	http.HandleFunc("/uploadHandler", uploadHandler)
	http.HandleFunc("/downloadActivity", downloadHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/removeActivity", removeHandler)
	http.HandleFunc("/editActivity", editHandler)
	http.HandleFunc("/searchCommentHandler", searchCommentHandler)
	http.Handle("/", http.FileServer(http.Dir("./Frontend")))

	//http.ListenAndServe(":80", nil)
	log.Fatalln(http.ListenAndServeTLS(":443", "Backend/main/cert.pem", "Backend/main/key.pem", nil))
}

type FrontendInfos struct {
	Activities map[int]Activity
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	oldCookie, err := r.Cookie("auth")
	if err == nil {
		delSessionKey(oldCookie.Value)
		cookie := http.Cookie{Name: "auth", Value: ""}
		http.SetCookie(w, &cookie)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	//Überprüfen ob die Sitzung des Nutzers noch gültig ist
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false { //Wenn nicht gültig, zurück zum Login
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else { //Wenn gültig übermittlete Form auswerten
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		activityIDString := r.Form.Get("actID") //Auslesen der geänderten Werte
		comment := r.Form.Get("comment")
		activityArt := r.Form.Get("actArt")
		activityID, err := strconv.Atoi(activityIDString)
		fmt.Println(err)
		if err == nil {
			editetAct := Activity{
				ActID:       activityID,
				Comment:     comment,
				UserID:      getUID(cookie.Value),
				Activityart: activityArt,
			}
			editActivity(editetAct) //Funktion die die bearbeiteten Werte in der CSV Datei ändert
		}
	}
}
func removeHandler(w http.ResponseWriter, r *http.Request) {
	//Überprüfen ob die Sitzung des Nutzers noch gültig ist
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false { //Wenn nicht gültig, zurück zum Login
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		activityIDString := r.Form.Get("actID") //Auslesen der zu löschenden Aktivität durch ID
		activityID, err := strconv.Atoi(activityIDString)
		fmt.Println(err)
		if err == nil {
			removeActivity(getUID(cookie.Value), activityID) //Funktion die die Zeile in der CSV Datei löscht
		}
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	//Überprüfen ob die Sitzung des Nutzers noch gültig ist
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false { //Wenn nicht gültig, zurück zum Login
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		activityIDString := r.Form.Get("actID")
		activityID, err := strconv.Atoi(activityIDString)
		readAcivityDB()              //Funktion zum aktualisieren der activityMap
		for k := range activityMap { //aktuelle acivityMap nach der zu donwloadenden Datei durchsuchen
			if activityMap[k].ActID == activityID {
				file := activityMap[k].Filename             //Filename/Filepfad auslesen
				downloadBytes, err := ioutil.ReadFile(file) //in Bytes zur Übermittlung ans Frontend packen

				if err != nil {
					fmt.Println(err)

				}

				mime := http.DetectContentType(downloadBytes) //Übermittlung ans Frontend

				fileSize := len(string(downloadBytes))

				// Festlegen der ResponseWriter
				w.Header().Set("Content-Type", mime)
				w.Header().Set("Content-Disposition", "attachment; filename="+file+"")
				w.Header().Set("Content-Length", strconv.Itoa(fileSize))
				// force it down the client's.....
				http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))
			}
		}
	}
}

//Funktion zum Aktualisieren der Aktivitäten auf dem Frontend
func viewDashboardHandler(w http.ResponseWriter, r *http.Request) {
	//Überprüfen ob die Sitzung des Nutzers noch gültig ist
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false { //Wenn nicht gültig, zurück zum Login
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		tmpl, error := template.ParseFiles("Frontend/dashboardTemplate.html")
		fmt.Println(error)

		var dataToTemplate = FrontendInfos{
			Activities: getDataForUser(getUID(cookie.Value)),
		}
		tmpl.Execute(w, dataToTemplate)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	pass := r.Form.Get("password")
	confirmPass := r.Form.Get("confirmPassword")
	userName := r.Form.Get("username")
	email := r.Form.Get("email")
	check, status := register(userName, email, pass, confirmPass)
	if check {
		cookie := http.Cookie{Name: "auth", Value: status}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		tmpl, error := template.ParseFiles("Frontend/RegisterTemplate.html")
		fmt.Println(error)
		tmpl.Execute(w, status)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	pass := r.Form.Get("password")
	userName := r.Form.Get("username")
	check, status := login(userName, pass)
	if check {
		cookie := http.Cookie{Name: "auth", Value: status}
		http.SetCookie(w, &cookie)
		fmt.Println("nice")
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		/*fmt.Println(status)
		http.Error(w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)*/
		tmpl, error := template.ParseFiles("Frontend/LoginTemplate.html")
		fmt.Println(error)
		tmpl.Execute(w, status)
	}
	//http.Redirect(w, r, "/home", http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	//Überprüfen ob die Sitzung des Nutzers noch gültig ist
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false { //Wenn nicht gültig, zurück zum Login
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {

		r.ParseMultipartForm(10 << 20)            //Auswertung der Form des Frontends
		file, handler, err := r.FormFile("datei") //Überprüfen ob der Nutzer eine Datei hochgeladen hat
		if err != nil {
			fmt.Println(file, "Error Retrieving the File")
			fmt.Println(err, handler)
			return
		}
		defer file.Close()
		if strings.HasSuffix(handler.Filename, ".gpx") { //Überprüfen ob es sich um eine GPX Datei handelt
			activity := r.FormValue("activity")
			kommentare := r.FormValue("kommentare")

			//Datei erstellen, in die die hochgeladene GPX Datei kopiert wird
			tempFile, err := ioutil.TempFile("DataStorage/GPX_Files", "gpxDatei*.gpx")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file) //Hochgeladene GPX Datei in Bytes packen
			if err != nil {
				fmt.Println(err)
			}
			tempFile.Write(fileBytes)                                               //Bytes in erstellete Datei laden, welche im DataStorage gespeichert wird
			uploadfile(tempFile.Name(), activity, kommentare, getUID(cookie.Value)) //Funktion die die hochgeladene Datei auswertet

		} else if strings.HasSuffix(handler.Filename, ".zip") { //Überprüfen ob es sich um eine ZIP Datei handelt
			//Temporäre Zip Datei erstellen, in die die hochgeladene Zip Datei kopiert wird
			tempFile, err := ioutil.TempFile("DataStorage/ZIP_Files", "zipDatei*.zip")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}
			tempFile.Write(fileBytes) //Bytes in erstellete Datei laden, welche im DataStorage gespeichert wird

			//Funktion zum Entpacken der Datei
			Unzip(tempFile.Name(), getUID(cookie.Value), r.FormValue("activity"), r.FormValue("kommentare"))
		} else {
			fmt.Println("No GPX or ZIP Data")
		}
		//Funktion zum Aktualisieren der Aktivitäten auf dem Frontend
		viewDashboardHandler(w, r)

	}
}
func Unzip(src string, uid int, actactivity string, komm string) {
	zipReader, _ := zip.OpenReader(src)
	//Zip Datei öffenen und Datei für Datei durchgehen
	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zippedFile.Close()
		if file.FileInfo().IsDir() { //Der erste Durchlauf ist nur das Verzeichnis
			log.Println("isDir")
		} else {
			if strings.HasSuffix(file.Name, ".gpx") { //Überprüfen ob es sich um eine GPX Datei handelt
				//Datei erstellen, in die die hochgeladene GPX Datei kopiert wird
				tempFile, err := ioutil.TempFile("DataStorage/GPX_Files", "gpxDatei*.gpx")
				if err != nil {
					fmt.Println(err)
				}
				defer tempFile.Close()
				filepath := tempFile.Name()
				outputFile, err := os.OpenFile(
					filepath,
					os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
					file.Mode(),
				)
				if err != nil {
					log.Fatal(err)
				}
				defer outputFile.Close()

				_, err = io.Copy(outputFile, zippedFile) //Daten in erstellete Datei kopieren, welche im DataStorage gespeichert wird
				if err != nil {
					log.Fatal(err)
				}
				uploadfile(tempFile.Name(), actactivity, komm, uid) //Funktion die die aktuelle GPX Datei auswertet
			}
		}
	}
}

func searchCommentHandler(w http.ResponseWriter, r *http.Request) {
	//Überprüfen ob die Sitzung des Nutzers noch gültig ist
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false { //Wenn nicht gültig, zurück zum Login
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		searchString := r.Form.Get("searchField")
		search(getUID(cookie.Value), searchString)
		var dataToTemplate = FrontendInfos{
			Activities: commentaryMap,
		}
		tmpl, error := template.ParseFiles("Frontend/dashboardTemplate.html")
		fmt.Println(error)
		tmpl.Execute(w, dataToTemplate)
	}
}
