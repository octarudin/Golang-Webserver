package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/iot")
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

func GetPZEMChartData(c *gin.Context) {
	rows, err := db.Query("SELECT timestamp, voltage, current, power FROM pzem ORDER BY timestamp ASC LIMIT 100")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()
	var data []map[string]interface{}
	for rows.Next() {
		var waktu string
		var voltage, current, power float64
		if err := rows.Scan(&waktu, &voltage, &current, &power); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		data = append(data, gin.H{
			"timestamp": waktu,
			"voltage":   voltage,
			"current":   current,
			"power":     power,
		})
	}
	c.JSON(200, data)
}

func GetXYMDChartData(c *gin.Context) {
	rows, err := db.Query("SELECT timestamp, temperature, humidity FROM xymd ORDER BY timestamp")
	if err != nil {
		log.Println("Query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var timestamp string
		var temperature, humidity float64
		if err := rows.Scan(&timestamp, &temperature, &humidity); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		data = append(data, gin.H{
			"timestamp":   timestamp,
			"temperature": temperature,
			"humidity":    humidity,
		})
	}

	c.JSON(http.StatusOK, data)
}
