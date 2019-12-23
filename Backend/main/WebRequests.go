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
	http.Handle("/", http.FileServer(http.Dir("./Frontend")))

	http.ListenAndServe(":80", nil)
}

func viewDashboardHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/", http.StatusFound)
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
	user := r.Form.Get("name")
	register(user, pass)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	pass := r.Form.Get("password")
	user := r.Form.Get("name")
	login(user, pass)
	http.Redirect(w, r, "/home", http.StatusFound)

}
