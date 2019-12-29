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

	//http.ListenAndServe(":80", nil)
	log.Fatalln(http.ListenAndServeTLS(":443", "Backend/main/cert.pem", "Backend/main/key.pem", nil))

}

func viewDashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		template, error := template.ParseFiles("Frontend/dashboardTemplate.html")
		fmt.Println(error)
		var lala = []Activity{
			Activity{
				actID:          1,
				UserID:         1,
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
				actID:          3,
				UserID:         4,
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
		}

		template.Execute(w, lala)
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
	/*cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
	*/
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
	tempFile.Write(fileBytes)
	uploadfile(tempFile.Name(), activity, kommentare)
	http.Redirect(w, r, "/landing.html", http.StatusFound)
	//}
}
