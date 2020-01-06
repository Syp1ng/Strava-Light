package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

type Gpx struct { //Benötigter Aufbau einer GPX Datei
	XMLName   xml.Name `xml:"gpx"` //Nicht benötigte Struktur wurde weggelassen
	Version   string   `xml:"version,attr"`
	Creator   string   `xml:"creator,attr"`
	Timestamp string   `xml:"metadata>time,omitempty"`
	Tracks    []GpxTrk `xml:"trk"`
}
type GpxTrk struct {
	XMLName  xml.Name    `xml:"trk"`
	Title    string      `xml:"name,omitempty"`
	Segments []GpxTrkSeg `xml:"trkseg"`
}
type GpxTrkSeg struct {
	XMLName xml.Name   `xml:"trkseg"`
	Points  []GpxPoint `xml:"trkpt"`
}
type GpxPoint struct {
	Lat       float64 `xml:"lat,attr"`
	Lon       float64 `xml:"lon,attr"`
	Timestamp string  `xml:"time"`
}

var gpxDocuments []*Gpx

func parseDoc(act Activity) Activity { //Funktion, die die hochgeladene Datei auswertet.

	gpxDoc, err := parseFile(act.Filename) //Funktion zum Laden der Datei in die Struktur
	if err != nil || gpxDoc == nil {

	}
	var actspeed = 0.0 //Definieren von Werten, die während der Auswertung immer wieder verändert werden
	var globalcounter = 0
	var kmCounter = 0.0
	var actkmSpeed = 0.0
	var timebetween = 0.0
	//var speedKM [100]float64
	var aktKM = 1000 //1000M -> 1KM

	//Drei Schleifen, weil die GPX Datei aus Track, Segment und Poiunts aufgebaut ist
	//und die Points zur Auswertung benötigt werden
	for i := 0; i < len(gpxDoc.Tracks); i++ {
		for j := 0; j < len(gpxDoc.Tracks[i].Segments); j++ {
			for k := 0; k < (len(gpxDoc.Tracks[i].Segments[j].Points))-1; k++ {
				distance2Points := distanceBetweenTwoPoints( //Funktion zur Berechnung der Entfernung zweier GPS Koordinaten
					gpxDoc.Tracks[i].Segments[j].Points[k].Lat, //Ausgelesen und Übergeben werden jeweils die
					gpxDoc.Tracks[i].Segments[j].Points[k].Lon, //Lat und Lon der beiden Punkte
					gpxDoc.Tracks[i].Segments[j].Points[k+1].Lat,
					gpxDoc.Tracks[i].Segments[j].Points[k+1].Lon)

				act.Distance += distance2Points
				if act.Distance >= float64(aktKM) { //Auswertung der Geschwindigkeit jedes Kilometers
					speedKM := actkmSpeed / kmCounter
					if act.AvgSpeedFastMS < speedKM { //Wenn die Durchschnittsseschwindigkeit dises Kilometers höher ist
						act.AvgSpeedFastKM = aktKM / 1000 //als die des bisher höchsten, soll diese überschrieben werden
						act.AvgSpeedFastMS = speedKM      //Zusätzlich wird der aktuelle Kilometer gespeichert
					}
					if act.AvgSpeedSlowMS > speedKM { //Wenn die Durchschnittsseschwindigkeit dises Kilometers geringer ist
						act.AvgSpeedSlowKM = aktKM / 1000 //als die des bisher geringsten, soll diese überschrieben werden
						act.AvgSpeedSlowMS = speedKM      //Zusätzlich wird der aktuelle Kilometer gespeichert
					}
					aktKM += 1000    //Absolvierte Kilometer werden erhöht um 1000M
					actkmSpeed = 0.0 //Zurücksetzten des aktuellen Kilometers um neue Berechung für neuen KM durchzuführen
					kmCounter = 0.0
				}
				//Überprüfung ob die GPX Datei einen Zeitstempel hat
				//Falls ja soll die Geschwindigkeit ausgerechnet werden
				if len(gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp) != 0 && len(gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp) != 0 { // noch Fehler

					//Funktion zur Berechnung der Geschwindikeit aus Entfernung der beiden Punkte und der dafür benötigten Zeit
					actspeed, timebetween = speed(distance2Points, gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp, gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp)
					act.Timestamp = gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp

					if actspeed > 0.5 { //Durchschnittsgeschwindigkeit nur auswerten, wenn Nutzer nicht steht
						act.Avgspeed = act.Avgspeed + actspeed //Addition der Aktuellen Geschwinbdigkeit für spätere
						actkmSpeed = actkmSpeed + actspeed     //Berechnung der Durchschnittsgeschwindigkeit
						globalcounter += 1
						kmCounter += 1
					} else { //Berechnung der Standzeit
						act.Standzeit = act.Standzeit + timebetween
					}
					if actspeed > act.HighSpeed { //Auswertung der Maximalgeschwindigkeit
						act.HighSpeed = actspeed
						act.Highspeedtime = gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp
					}
				} else {
					act.Timestamp = "NoTime"
					log.Println("NoTime")
				}
			}
		}
	}
	act.Avgspeed = act.Avgspeed / float64(globalcounter) //Berechnung der Durchschnittsgeschwindigkeit
	act.Distance = math.Round(act.Distance/10) / 100     //km + Rundung auf 2 Nachkommastellen
	act.HighSpeed = math.Round(act.HighSpeed*100) / 100
	return act
}

//Funktion zur Berechnung der Geschwindikeit aus Entfernung der beiden Punkte und der dafür benötigten Zeit
func speed(distance float64, timestamp1 string, timestamp2 string) (float64, float64) {
	var timebetween = 0.0
	layout := "2006-01-02T15:04:05.000Z" //Auswertung der Zeitstempel
	t, err := time.Parse(layout, timestamp1)
	t1, err := time.Parse(layout, timestamp2)
	if err != nil {
		layout := "2006-01-02T15:04:05Z"
		t, err = time.Parse(layout, timestamp1)
		t1, err = time.Parse(layout, timestamp2)
		if err != nil {
			fmt.Println(err)
		}
	}
	timebetween = float64(t1.Sub(t)) / 1000000000 //Berechnung der vergangenen Zeit zwischen den beiden Zeitstempel
	return distance / timebetween, timebetween    //Berechnung der Geschwindigkeit

}

//Funktion zum Laden der benötigten Daten aus der GPX Datei
func parseFile(filename string) (*Gpx, error) {
	gpxFile, err := os.Open(filename) //Filename zum finden der Datei im DataStorage/GPX_Files
	if err != nil {
		return nil, err
	}

	defer gpxFile.Close()
	b, err := ioutil.ReadAll(gpxFile) //Auslesen der Datei
	if err != nil {
		return nil, err
	}

	g := &Gpx{}
	xml.Unmarshal(b, &g) //Benötigte Dateien der Datei in das GPX Struck schreiben
	return g, nil        //GPX Struck zur Auswertung zurückgeben
}

func ToRad(x float64) float64 {
	return x / 180. * math.Pi
}

//Funktion zur Berechnung der Entfernung zweier GPS Koordinaten
//http://www.movable-type.co.uk/scripts/latlong.html
func distanceBetweenTwoPoints(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := ToRad(lat1 - lat2)
	dLon := ToRad(lon1 - lon2)
	thisLat1 := ToRad(lat1)
	thisLat2 := ToRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := 6371000 * c //Erdradius

	return d
}
