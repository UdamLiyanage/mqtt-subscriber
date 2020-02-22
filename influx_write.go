package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
	"io"
	"log"
	"os"
	"strconv"
)

var influx client.Client

func writeData(msg []byte) {
	println(string(msg))
	r := csv.NewReader(bytes.NewBuffer(msg))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print("Error: ", err)
		}
		writeToInflux(record)
	}
}

func writeToInflux(rec []string) {
	if len(rec) == 2 {
		phDevice(rec)
	} else if len(rec) == 4 {
		conditionsDevice(rec)
	} else {
		log.Println("Wrong Format")
	}
}

func getMean() string {
	q := fmt.Sprintf("SELECT MEAN(%s) FROM %s WHERE time > now() - 5m LIMIT 12", "temp", "conditions")
	query := client.Query{
		Command:  q,
		Database: os.Getenv("INFLUX_DB_DATABASE"),
	}
	resp, err := influx.Query(query)
	if err != nil {
		log.Println("Error: ", err)
	}
	res := resp.Results[0].Series[0].Values[0][1].(json.Number).String()
	return res
}

func connect() {
	var err error
	influx, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: os.Getenv("INFLUX_DB"),
	})
	if err != nil {
		log.Println("Error: ", err)
	}
}

func ParseFloat(s string, bitSize int) float64 {
	res, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		log.Println("Error: ", err)
	}
	return res
}

func phDevice(rec []string) {
	point, err := client.NewPoint(
		"ph",
		map[string]string{
			"serial_number": rec[1],
		},
		map[string]interface{}{
			"ph": ParseFloat(rec[0], 32),
		})
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: os.Getenv("INFLUX_DB_DATABASE"),
	})
	if err != nil {
		log.Println("Error: ", err)
	}

	bps.AddPoint(point)

	err = influx.Write(bps)
	if err != nil {
		log.Println("Error: ", err)
	}
}

func conditionsDevice(rec []string) {
	point, err := client.NewPoint(
		"conditions",
		map[string]string{
			"serial_number": rec[3],
		},
		map[string]interface{}{
			"temp":     ParseFloat(rec[0], 32),
			"humidity": ParseFloat(rec[1], 32),
			"light":    ParseFloat(rec[2], 32),
		})
	if err != nil {
		log.Println("Error: ", err)
	}
	bps, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: os.Getenv("INFLUX_DB_DATABASE"),
	})
	if err != nil {
		log.Println("Error: ", err)
	}

	bps.AddPoint(point)

	err = influx.Write(bps)
	if err != nil {
		log.Println("Error: ", err)
	}
}
