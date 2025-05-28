package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	StartMQTT()

	router := gin.Default()

	// ✅ Register API route dulu
	router.GET("/api/pzem", GetPZEMData)
	router.GET("/api/xymd", GetXYMDData)
	router.GET("/api/pzem/chart", GetPZEMChartData)
	router.GET("/api/xymd/chart", GetXYMDChartData)

	// ✅ Baru static file di root
	router.Static("/web", "./frontend")

	// webserver berhasil di run,, jangan lupa
	// go run ./backend/
	// http://localhost:8080/web/
	router.Run(":8080")
}
