package main

import (
	"archive/zip"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func SetupLinks() {
	http.HandleFunc("/home", viewDashboardHandler)
	http.HandleFunc("/registrationHandler", registerHandler)
	http.HandleFunc("/loginHandler", loginHandler)
	http.HandleFunc("/uploadHandler", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir("./Frontend")))
	http.HandleFunc("/logout", logoutHandler)

	//http.ListenAndServe(":80", nil)
	log.Fatalln(http.ListenAndServeTLS(":443", "Backend/main/cert.pem", "Backend/main/key.pem", nil))
}

type FrontendInf struct {
	Activities []Activity
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

func viewDashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {

		//hier geht net scheißdreck
		tmpl, error := template.ParseFiles("Frontend/dashboardTemplate.html")
		fmt.Println(error)

		var lala = FrontendInf{
			Activities: getDataForUser(getUID(cookie.Value)),
		}
		tmpl.Execute(w, lala)
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
		fmt.Println(status)
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
		fmt.Println(status)
		http.Error(w,
			http.StatusText(http.StatusUnauthorized),
			http.StatusUnauthorized)
	}

	//http.Redirect(w, r, "/home", http.StatusFound)

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {

		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("datei")
		if err != nil {
			fmt.Println(file, "Error Retrieving the File")
			fmt.Println(err, handler)
			return
		}
		defer file.Close()
		if strings.HasSuffix(handler.Filename, ".gpx") {
			activity := r.FormValue("activity")
			kommentare := r.FormValue("kommentare")
			tempFile, err := ioutil.TempFile("DataStorage/GPX_Files", "gpxDatei*.gpx")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}
			tempFile.Write(fileBytes)
			uploadfile(tempFile.Name(), activity, kommentare, getUID(cookie.Value))

		} else if strings.HasSuffix(handler.Filename, ".zip") {
			tempFile, err := ioutil.TempFile("DataStorage/ZIP_Files", "zipDatei*.zip")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}
			tempFile.Write(fileBytes)
			gpxFiles, err := Unzip(tempFile.Name(), getUID(cookie.Value), r.FormValue("activity"), r.FormValue("kommentare"))
			if err != nil {
				log.Fatal(err)
			}
			log.Println(gpxFiles)
		} else {
			fmt.Println("No GPX or ZIP Data")
		}
		viewDashboardHandler(w, r)

	}
}
func Unzip(src string, uid int, actactivity string, komm string) ([]string, error) {
	zipReader, _ := zip.OpenReader(src)
	for _, file := range zipReader.Reader.File {
		zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zippedFile.Close()
		if file.FileInfo().IsDir() {
			log.Println("isDir")
		} else {
			if strings.HasSuffix(file.Name, ".gpx") {
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

				_, err = io.Copy(outputFile, zippedFile)
				if err != nil {
					log.Fatal(err)
				}
				uploadfile(tempFile.Name(), actactivity, komm, uid)
			}
		}
	}

	return nil, nil
}
