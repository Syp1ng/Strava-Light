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

//http://www.movable-type.co.uk/scripts/latlong.html

func main() {
	parseDoc()
}

type Gpx struct {
	XMLName   xml.Name `xml:"gpx"`
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

func parseDoc() {

	gpxDoc, err := parseFile()
	if err != nil || gpxDoc == nil {

	}
	var distance = 0.0
	var speed = 0.0
	var highSpeed = 0.0
	var highSpeedTime = ""
	var timebetween = 0.0
	for i := 0; i < len(gpxDoc.Tracks); i++ {
		for j := 0; j < len(gpxDoc.Tracks[i].Segments); j++ {
			for k := 0; k < (len(gpxDoc.Tracks[i].Segments[j].Points))-1; k++ {
				//log.Print(gpxDoc.Tracks[i].Segments[j].Points[k].Lat , " ", gpxDoc.Tracks[i].Segments[j].Points[k].Lon)
				distance2Points := distanceBetweenTwoPoints(
					gpxDoc.Tracks[i].Segments[j].Points[k].Lat,
					gpxDoc.Tracks[i].Segments[j].Points[k].Lon,
					gpxDoc.Tracks[i].Segments[j].Points[k+1].Lat,
					gpxDoc.Tracks[i].Segments[j].Points[k+1].Lon)

				if gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp != "" || gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp != "" { // noch Fehler
					layout := "2006-01-02T15:04:05.000Z"
					str := gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp
					str1 := gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp
					t, err := time.Parse(layout, str)
					t1, err := time.Parse(layout, str1)
					if err != nil {
						fmt.Println(err)
					}
					timebetween = float64(t1.Sub(t)) / 1000000000
					speed = distance2Points / timebetween
					if speed > highSpeed {
						highSpeed = speed
						highSpeedTime = gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp
					}
				}
				distance += distance2Points
				log.Print(distance2Points, timebetween, speed, highSpeed, highSpeedTime)

			}
		}
	}
}

func parseFile() (*Gpx, error) {
	gpxFile, err := os.Open("GPX_Files/2019-09-14_15-14.gpx")
	if err != nil {
		return nil, err
	}

	defer gpxFile.Close()
	b, err := ioutil.ReadAll(gpxFile)
	if err != nil {
		return nil, err
	}

	g := &Gpx{}
	xml.Unmarshal(b, &g)
	return g, nil
}

func ToRad(x float64) float64 {
	return x / 180. * math.Pi
}

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
