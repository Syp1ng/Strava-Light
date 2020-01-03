package main

import (
	//"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestUploadFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile("DataStorage/GPX_Files", "gpxDatei*.gpx")
	uploadfile(tempFile.Name(), "Laufen", "Kommentar", 1)
	//assert.
}
