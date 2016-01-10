package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"strconv"
	"math"
)

const (
// According to Wikipedia, the Earth's radius is about 6,371km
	EARTH_RADIUS = 6371
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func value(v interface{}, e error) interface{} {
	check(e)
	return v
}

type WayPointData struct {
	Latitude  string `xml:"lat,attr"`
	Longitude string `xml:"lon,attr"`
	Elevation string `xml:"ele"`
	Time string `xml:"time"`
}

func (wpd *WayPointData) wayPoint() wayPoint {
	var wp wayPoint

	wp.time = value(time.Parse("2006-01-02T15:04:05.000Z", wpd.Time)).(time.Time)
	wp.latitude = value(strconv.ParseFloat(wpd.Latitude, 64)).(float64)
	wp.longitude = value(strconv.ParseFloat(wpd.Longitude, 64)).(float64)
	wp.elevation = value(strconv.ParseFloat(wpd.Elevation, 64)).(float64)

	return wp
}

type point struct {
	latitude float64
	longitude float64
	elevation float64
}

func (p *point) isZero() bool {
	return p.latitude == 0 && p.longitude == 0 && p.elevation == 0
}

// Calculates the Haversine distance between two points in kilometers.
// Copied from https://github.com/kellydunn/golang-geo/blob/master/point.go
func (this *point) distanceTo(that point) float64 {
	dLat := (that.latitude - this.latitude) * (math.Pi / 180.0)
	dLon := (that.longitude - this.longitude) * (math.Pi / 180.0)

	lat1 := this.latitude * (math.Pi / 180.0)
	lat2 := that.latitude * (math.Pi / 180.0)

	a1 := math.Sin(dLat / 2) * math.Sin(dLat / 2)
	a2 := math.Sin(dLon/2) * math.Sin(dLon/2) * math.Cos(lat1) * math.Cos(lat2)

	a := a1 + a2

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTH_RADIUS * c
}

type wayPoint struct {
	time time.Time
	point
}

type TrackInfo struct {
	filePath  string
	startTime time.Time
	duration  time.Duration
	distance  float64
	ascent    int
	descent   int
	lastPoint point
}

// FilePath returns the path of the file from which GPX summary has been extracted.
func (info *TrackInfo) FilePath() string { return info.filePath }

// StartTime returns the time at which the track starts.
func (info *TrackInfo) StartTime() time.Time { return info.startTime }

// Duration returns the total duration of the track.
func (info *TrackInfo) Duration() time.Duration { return info.duration }

// Distance returns the total distance of the track in kilometers.
func (info *TrackInfo) Distance() float64 { return info.distance }

// Ascent returns the cumulative ascent of the track in meters.
func (info *TrackInfo) Ascent() int { return info.ascent }

// Descent returns the cumulative descent of the track in meters.
func (info *TrackInfo) Descent() int { return info.descent }

// Speed returns the average speed for the track in kilometers/hour.
func (info *TrackInfo) Speed() float64 { return info.distance / info.duration.Hours() }

// Pace returns the average pace for the trace in minutes/kilometer.
func (info *TrackInfo) Pace() float64 { return info.duration.Minutes() / info.distance }

func (info *TrackInfo) append(wp wayPoint) {
	if info.startTime.IsZero() {
		info.startTime = wp.time
	}

	info.duration = wp.time.Sub(info.startTime)

	if !info.lastPoint.isZero() {
		info.distance += info.lastPoint.distanceTo(wp.point)

		if info.lastPoint.elevation != 0 {
			eleDiff := math.Abs(wp.elevation - info.lastPoint.elevation)
			if info.lastPoint.elevation < wp.elevation {
				info.ascent += int(eleDiff)
			} else {
				info.descent += int(eleDiff)
			}
		}
	}

	info.lastPoint = wp.point
}

func (info *TrackInfo) Format() string {
	return fmt.Sprintf("File: %s\nTime: %s\nDuration: %s\nDistance: %.1fkm\nAscent: %dm\nDescent: %dm\nSpeed: %.1fkm/h\nPace: %.1fmin/km",
		info.FilePath(), info.StartTime().Local().Format(time.RFC1123), info.Duration(),
		info.Distance(), info.Ascent(), info.Descent(),
		info.Speed(),
		info.Pace())
}

// Process extracts the GPX track summary from the given file.
func Process(filePath string) TrackInfo {
	file := value(os.Open(filePath)).(*os.File)
	defer file.Close()

	info := TrackInfo{filePath: filePath}
	decoder := xml.NewDecoder(file)

	for {
		t, err := decoder.Token()

		if err == io.EOF {
			break
		}
		check(err)

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "trkpt" {
				var wayPoint WayPointData

				decoder.DecodeElement(&wayPoint, &se)

				info.append(wayPoint.wayPoint())
			}
		}
	}

	return info
}

func main() {
	flag.Parse()

	for _, filePath := range flag.Args() {
		info := Process(filePath)
		fmt.Println(info.Format())
	}
}
