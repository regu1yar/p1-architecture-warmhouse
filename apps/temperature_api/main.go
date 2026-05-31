package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort = "8081"
	unit        = "celsius"
	status      = "active"
	sensorType  = "temperature"
	minTemp     = -50.0
	maxTemp     = 40.0
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

var sensorIDToLocation = map[string]string{
	"1": "Living Room",
	"2": "Bedroom",
	"3": "Kitchen",
}

var locationToSensorID = map[string]string{
	"Living Room": "1",
	"Bedroom":     "2",
	"Kitchen":     "3",
}

func resolveLocation(location, sensorID string) (string, string) {
	if location == "" {
		if l, ok := sensorIDToLocation[sensorID]; ok {
			location = l
		} else {
			location = "Unknown"
		}
	}
	if sensorID == "" {
		if id, ok := locationToSensorID[location]; ok {
			sensorID = id
		} else {
			sensorID = "0"
		}
	}
	return location, sensorID
}

func randTemp() float64 {
	return minTemp + rand.Float64()*(maxTemp-minTemp)
}

func buildResponse(location, sensorID string) TemperatureResponse {
	location, sensorID = resolveLocation(location, sensorID)
	return TemperatureResponse{
		Value:       randTemp(),
		Unit:        unit,
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      status,
		SensorID:    sensorID,
		SensorType:  sensorType,
		Description: "Temperature in " + location,
	}
}

func getTemperatureByLocation(c *gin.Context) {
	c.JSON(http.StatusOK, buildResponse(c.Query("location"), ""))
}

func getTemperatureByID(c *gin.Context) {
	c.JSON(http.StatusOK, buildResponse("", c.Param("sensorID")))
}

func listenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if strings.HasPrefix(port, ":") {
		return port
	}
	return ":" + port
}

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/temperature", getTemperatureByLocation)
	r.GET("/temperature/:sensorID", getTemperatureByID)

	addr := listenAddr()
	log.Printf("temperature-api listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
