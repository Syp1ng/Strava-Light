package main

import (
	"fmt"
	"log"
	"net/http"
)

func SetupLinks() {
	http.HandleFunc("/home", viewDashboardHandler)
	http.HandleFunc("/registrationHandler", registerHandler)
	http.HandleFunc("/loginHandler", loginHandler)
	http.HandleFunc("/uploadHandler", uploadHandler)
	http.Handle("/", http.FileServer(http.Dir("./Frontend")))

	http.ListenAndServe(":80", nil)
}

func viewDashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		//get userData.....
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	pass := r.Form.Get("password")
	userName := r.Form.Get("username")
	email := r.Form.Get("email")
	check, status := register(email, userName, pass)
	if check {
		cookie := http.Cookie{Name: "auth", Value: status}
		http.SetCookie(w, &cookie)

	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	pass := r.Form.Get("password")
	email := r.Form.Get("email")
	check, status := login(email, pass)
	if check {
		cookie := http.Cookie{Name: "auth", Value: status}
		http.SetCookie(w, &cookie)
		fmt.Println("nice")
	}
	//http.Redirect(w, r, "/home", http.StatusFound)

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		//get userData.....
	}
}
