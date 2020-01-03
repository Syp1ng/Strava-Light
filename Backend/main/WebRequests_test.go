package main

import (
	_ "github.com/stretchr/testify"
	_ "testing"
)

//////NOCH nichts ausser Signatur
/*
func TestSetupLinks() {
	http.HandleFunc("/home", viewDashboardHandler)
	http.HandleFunc("/registrationHandler", registerHandler)
	http.HandleFunc("/loginHandler", loginHandler)
	http.HandleFunc("/uploadHandler", uploadHandler)
	http.HandleFunc("/downloadActivity", downloadHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/removeActivity", removeHandler)
	http.HandleFunc("/editActivity", editHandler)
	http.Handle("/", http.FileServer(http.Dir("./Frontend")))

	//http.ListenAndServe(":80", nil)
	log.Fatalln(http.ListenAndServeTLS(":443", "Backend/main/cert.pem", "Backend/main/key.pem", nil))
}

func TestLogoutHandler(t *testing.T) {
	oldCookie, err := r.Cookie("auth")
	if err == nil {
		delSessionKey(oldCookie.Value)
		cookie := http.Cookie{Name: "auth", Value: ""}
		http.SetCookie(w, &cookie)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func TestEditHandler(t *testing.T) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		activityIDString := r.Form.Get("actID")
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
			editActivity(editetAct)
		}
	}
}
func TestRemoveHandler(t *testing.T) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		activityIDString := r.Form.Get("actID")
		activityID, err := strconv.Atoi(activityIDString)
		fmt.Println(err)
		if err == nil {
			removeActivity(getUID(cookie.Value), activityID)
		}
	}
}

func TestDownloadHandler(t *testing.T) {

	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		activityIDString := r.Form.Get("actID")
		activityID, err := strconv.Atoi(activityIDString)
		readAcivityDB()
		for k := range activityMap {
			if activityMap[k].ActID == activityID {
				file := activityMap[k].Filename
				downloadBytes, err := ioutil.ReadFile(file)

				if err != nil {
					fmt.Println(err)
				}

				// set the default MIME type to send
				mime := http.DetectContentType(downloadBytes)

				fileSize := len(string(downloadBytes))

				// Generate the server headers
				w.Header().Set("Content-Type", mime)
				w.Header().Set("Content-Disposition", "attachment; filename="+file+"")
				w.Header().Set("Content-Length", strconv.Itoa(fileSize))
				// force it down the client's.....
				http.ServeContent(w, r, file, time.Now(), bytes.NewReader(downloadBytes))

			}

		}

	}

}

func TestViewDashboardHandler(t *testing.T) {
	cookie, err := r.Cookie("auth")
	if err != nil || checkSessionKey(cookie.Value) == false {
		fmt.Printf("No cookie or invalid Session")
		http.Redirect(w, r, "/Login.html", http.StatusFound)
	} else {
		tmpl, error := template.ParseFiles("Frontend/dashboardTemplate.html")
		fmt.Println(error)
		var dataToTemplate = FrontendInf{
			Activities: getDataForUser(getUID(cookie.Value)),
		}
		tmpl.Execute(w, dataToTemplate)
	}
}

func TestRegisterHandler(t *testing.T) {
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

func TestLoginHandler(t *testing.T) {
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
		tmpl, error := template.ParseFiles("Frontend/LoginTemplate.html")
		fmt.Println(error)
		tmpl.Execute(w, status)
	}

	//http.Redirect(w, r, "/home", http.StatusFound)

}

func TestUploadHandler(t *testing.T) {
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
func TestUnzip(t *testing.T) {
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
}*/
