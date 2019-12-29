package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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
		//buffer,err := file.Read(make([]byte, 512))
		if err != nil {

		}
		tempFile.Write(fileBytes)
		uploadfile(tempFile.Name(), activity, kommentare, getUID(cookie.Value))
		viewDashboardHandler(w, r)
		//http.Redirect(w, r, "/dashboardTemplate.html", http.StatusFound)
	}
}
