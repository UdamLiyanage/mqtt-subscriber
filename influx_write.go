package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"net/url"
	"strconv"
	_ "time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client"
)

var influx *client.Client

func writeData(msg []byte) {
	println(string(msg))
	connect()
	r := csv.NewReader(bytes.NewBuffer(msg))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		writeToInflux(record)
	}
}

func writeToInflux(rec []string) {
	if len(rec) == 2 {
		phDevice(rec)
	} else if len(rec) == 4 {
		conditionsDevice(rec)
	}
}

func connect() {
	u, err := url.Parse("http://54.67.77.75:8086")
	if err != nil {
		panic(err)
	}
	influx, err = client.NewClient(client.Config{URL: *u})
	if err != nil {
		panic(err)
	}
}

func ParseFloat(s string, bitSize int) float64 {
	res, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		panic(err)
	}
	return res
}

func phDevice(rec []string) {
	var pts = make([]client.Point, 2)
	pts[0] = client.Point{
		Measurement: "ph",
		Tags: map[string]string{
			"serial_number": rec[1],
		},
		Fields: map[string]interface{}{
			"ph": ParseFloat(rec[0], 32),
		},
	}
	bps := client.BatchPoints{
		Points:   pts,
		Database: "button_test",
	}

	_, err := influx.Write(bps)
	if err != nil {
		panic(err)
	}
}

func conditionsDevice(rec []string) {
	var pts = make([]client.Point, 2)
	pts[0] = client.Point{
		Measurement: "conditions",
		Tags: map[string]string{
			"serial_number": rec[3],
		},
		Fields: map[string]interface{}{
			"temp":     ParseFloat(rec[0], 32),
			"humidity": ParseFloat(rec[1], 32),
			"light":    ParseFloat(rec[2], 32),
		},
	}
	bps := client.BatchPoints{
		Points:   pts,
		Database: "button_test",
	}

	_, err := influx.Write(bps)
	if err != nil {
		panic(err)
	}
}
