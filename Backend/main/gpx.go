package main

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"math"
	"os"
)

//http://www.movable-type.co.uk/scripts/latlong.html

func main() {
	parseDoc()
}

type Gpx struct {
	XMLName   xml.Name `xml:"gpx"`
	Version   string   `xml:"version"`
	Creator   string   `xml:"creator"`
	Title     string   `xml:"metadata>text"`
	Timestamp string   `xml:"metadata>time"`
	Tracks    []GpxTrk `xml:"trk"`
}
type GpxTrk struct {
	XMLName  xml.Name    `xml:"trk"`
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
	log.Println(gpxDoc.XMLName)
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

func DistanceBetweenTwoPoints(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := ToRad(lat1 - lat2)
	dLon := ToRad(lon1 - lon2)
	thisLat1 := ToRad(lat1)
	thisLat2 := ToRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(thisLat1)*math.Cos(thisLat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := 6371000 * c //Erdradius

	return d
}
