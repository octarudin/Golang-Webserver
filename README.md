Project Template: Golang + MySQL + MQTT + Bootstrap
===================================================

📁 project-root/ 
├── backend/
│   ├── main.go
│   ├── mqtt_client.go
│   ├── db.go
│   └── models.go
├── frontend/
│   └── index.html
├── go.mod
└── README.md

--------------------------------------------------
1. backend/main.go
--------------------------------------------------
package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	StartMQTT()

	r := gin.Default()
	r.Static("/", "../frontend")
	r.GET("/api/pzem", GetPZEMData)
	r.GET("/api/xymd", GetXYMDData)

	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}

--------------------------------------------------
2. backend/mqtt_client.go
--------------------------------------------------
package main

import (
	"encoding/json"
	"log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PZEMPayload struct {
	Voltage     float64 `json:"voltage"`
	Current     float64 `json:"current"`
	Power       float64 `json:"power"`
	Energy      float64 `json:"energy"`
	Frequency   float64 `json:"frequency"`
	PowerFactor float64 `json:"power_factor"`
}

type XYMDPayload struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func StartMQTT() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("go_dashboard")
	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	client.Subscribe("esp32/pzem", 0, func(c mqtt.Client, m mqtt.Message) {
		var data PZEMPayload
		if err := json.Unmarshal(m.Payload(), &data); err == nil {
			InsertPZEMData(data)
		}
	})

	client.Subscribe("esp32/xymd", 0, func(c mqtt.Client, m mqtt.Message) {
		var data XYMDPayload
		if err := json.Unmarshal(m.Payload(), &data); err == nil {
			InsertXYMDData(data)
		}
	})
}

--------------------------------------------------
3. backend/db.go
--------------------------------------------------
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/iot")
	if err != nil {
		log.Fatal(err)
	}
}

func InsertPZEMData(data PZEMPayload) {
	db.Exec(`INSERT INTO pzem (voltage, current, power, energy, frequency, power_factor, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, NOW())`,
		data.Voltage, data.Current, data.Power, data.Energy, data.Frequency, data.PowerFactor)
}

func InsertXYMDData(data XYMDPayload) {
	db.Exec(`INSERT INTO xymd (temperature, humidity, timestamp)
		VALUES (?, ?, NOW())`,
		data.Temperature, data.Humidity)
}

func GetPZEMData(c *gin.Context) {
	rows, _ := db.Query("SELECT * FROM pzem ORDER BY timestamp DESC LIMIT 20")
	var result []map[string]interface{}
	for rows.Next() {
		var v, i, p, e, f, pf float64
		var t string
		rows.Scan(&v, &i, &p, &e, &f, &pf, &t)
		result = append(result, gin.H{"voltage": v, "current": i, "power": p, "energy": e, "frequency": f, "power_factor": pf, "timestamp": t})
	}
	c.JSON(200, result)
}

func GetXYMDData(c *gin.Context) {
	rows, _ := db.Query("SELECT * FROM xymd ORDER BY timestamp DESC LIMIT 20")
	var result []map[string]interface{}
	for rows.Next() {
		var temp, hum float64
		var t string
		rows.Scan(&temp, &hum, &t)
		result = append(result, gin.H{"temperature": temp, "humidity": hum, "timestamp": t})
	}
	c.JSON(200, result)
}

--------------------------------------------------
4. frontend/index.html
--------------------------------------------------
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>IoT Dashboard</title>
  <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body class="p-4">
  <div class="container">
    <h2>PZEM Data</h2>
    <table class="table table-bordered" id="pzem-table"></table>

    <h2>XY-MD02 Data</h2>
    <table class="table table-bordered" id="xymd-table"></table>
  </div>
  <script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
  <script>
    $.getJSON("/api/pzem", function(data) {
      let html = "<tr><th>Voltage</th><th>Current</th><th>Power</th><th>Energy</th><th>Frequency</th><th>PF</th><th>Time</th></tr>";
      data.forEach(d => {
        html += `<tr><td>${d.voltage}</td><td>${d.current}</td><td>${d.power}</td><td>${d.energy}</td><td>${d.frequency}</td><td>${d.power_factor}</td><td>${d.timestamp}</td></tr>`;
      });
      $("#pzem-table").html(html);
    });

    $.getJSON("/api/xymd", function(data) {
      let html = "<tr><th>Temperature</th><th>Humidity</th><th>Time</th></tr>";
      data.forEach(d => {
        html += `<tr><td>${d.temperature}</td><td>${d.humidity}</td><td>${d.timestamp}</td></tr>`;
      });
      $("#xymd-table").html(html);
    });
  </script>
</body>
</html>

--------------------------------------------------
5. go.mod
--------------------------------------------------
module iot-dashboard

go 1.20

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/gin-gonic/gin v1.9.1
	github.com/go-sql-driver/mysql v1.6.0

	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/gin-gonic/gin v1.10.1
	github.com/go-sql-driver/mysql v1.9.2
)

--------------------------------------------------
6. SQL Table Setup (MariaDB)
--------------------------------------------------
CREATE DATABASE iot;
USE iot;

CREATE TABLE pzem (
	voltage DOUBLE,
	current DOUBLE,
	power DOUBLE,
	energy DOUBLE,
	frequency DOUBLE,
	power_factor DOUBLE,
	timestamp DATETIME
);

CREATE TABLE xymd (
	temperature DOUBLE,
	humidity DOUBLE,
	timestamp DATETIME
);

--------------------------------------------------
7. Jalankan
--------------------------------------------------
go mod tidy
go run backend/main.go -- fails
go run ./backend/

Lalu akses: http://localhost:8080
