package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"
)

//http://www.movable-type.co.uk/scripts/latlong.html

//func main() {
//parseDoc("GPX_Files/2019-09-21_15-54.gpx")
//}

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

func parseDoc(filename string) {

	gpxDoc, err := parseFile(filename)
	if err != nil || gpxDoc == nil {

	}
	var distance = 0.0
	var highSpeed = 0.0
	var highSpeedTime = ""
	var actspeed = 0.0
	var avgspeed = 0.0
	var globalcounter = 0.0
	var kmCounter = 0.0
	var actkmSpeed = 0.0
	var timebetween = 0.0
	var standzeit = 0.0
	var speedKM [100]float64
	var aktKM = 1000
	for i := 0; i < len(gpxDoc.Tracks); i++ {
		for j := 0; j < len(gpxDoc.Tracks[i].Segments); j++ {
			for k := 0; k < (len(gpxDoc.Tracks[i].Segments[j].Points))-1; k++ {
				//log.Print(gpxDoc.Tracks[i].Segments[j].Points[k].Lat , " ", gpxDoc.Tracks[i].Segments[j].Points[k].Lon)
				distance2Points := distanceBetweenTwoPoints( //6.1
					gpxDoc.Tracks[i].Segments[j].Points[k].Lat,
					gpxDoc.Tracks[i].Segments[j].Points[k].Lon,
					gpxDoc.Tracks[i].Segments[j].Points[k+1].Lat,
					gpxDoc.Tracks[i].Segments[j].Points[k+1].Lon)
				distance += distance2Points
				if distance >= float64(aktKM) {
					fmt.Println(actkmSpeed, kmCounter)
					speedKM[aktKM/1000] = actkmSpeed / kmCounter //12
					aktKM += 1000
					actkmSpeed = 0.0
					kmCounter = 0.0
				}
				if gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp != "" || gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp != "" { // noch Fehler
					actspeed, timebetween = speed(distance2Points, gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp, gpxDoc.Tracks[i].Segments[j].Points[k+1].Timestamp)
					if actspeed > 0.5 { //6.4
						avgspeed = avgspeed + actspeed //6.2
						actkmSpeed = actkmSpeed + actspeed
						globalcounter += 1
						kmCounter += 1
					} else { //6.3
						standzeit = standzeit + timebetween
					}
					if actspeed > highSpeed {
						highSpeed = actspeed
						highSpeedTime = gpxDoc.Tracks[i].Segments[j].Points[k].Timestamp
					}
				}

				//log.Print(distance, distance2Points, actspeed, highSpeed, highSpeedTime)
			}
			avgspeed = avgspeed / globalcounter      //noch Fehler
			distance = math.Round(distance/10) / 100 //km + Rundung auf 2 Nachkommastellen
			highSpeed = math.Round(highSpeed*100) / 100
			fmt.Println(distance, highSpeed, highSpeedTime, avgspeed, globalcounter, standzeit, speedKM[4])
		}
	}
}
func speed(distance float64, timestamp1 string, timestamp2 string) (float64, float64) {
	var timebetween = 0.0
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, timestamp1)
	t1, err := time.Parse(layout, timestamp2)
	if err != nil {
		fmt.Println(err)
	}
	timebetween = float64(t1.Sub(t)) / 1000000000
	return distance / timebetween, timebetween

}

func parseFile(filename string) (*Gpx, error) {
	gpxFile, err := os.Open(filename)
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
