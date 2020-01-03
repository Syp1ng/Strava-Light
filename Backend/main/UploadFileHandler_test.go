package main

import (
	"encoding/csv"
	"log"
	"strconv"

	//"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestUploadFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile("DataStorage/GPX_Files", "gpxDatei*.gpx")
	uploadfile(tempFile.Name(), "Laufen", "Kommentar", 1)
	//assert.
	activityDataLenght, error := os.Open(dbLocationActivity)
	if error == nil {
		readerforlengh := csv.NewReader(activityDataLenght)
		for {
			lin, err := readerforlengh.Read()
			if err == nil {
				if lin[2] == //da den PFadlalala {
					//dann is in actdb drin
				}

			} else {
				break
				log.Println(lin)
			}
		}
	}
}
