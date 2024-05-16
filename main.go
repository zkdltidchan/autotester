package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zkdltidchan/autotester/handlers"
)

const (
	MaxBrokerCount           = 3
	RestaterCheckInterval    = 1 * time.Second
	BrokerReadyCheckInterval = 1 * time.Second
	// BaseDir                  = `C:\Users\SJSJ\Desktop\wowsan`
	// SeedBatchFile            = `scripts\broker.bat`
	// SimulatorFile = `cmd\simulator\main.go`
	BaseDir       = `/Users/zkdltid/Desktop/go-test2/pubsubtest`
	SeedBatchFile = `scripts/broker.sh`
	SimulatorFile = `cmd/simulator/main.go`
)

func main() {
	r := gin.Default()

	logger := logrus.New()
	logFile, err := os.OpenFile(fmt.Sprintf("server_%s.log", time.Now().Format("20060102_150405")), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}
	logger.SetOutput(logFile)

	service := handlers.NewServer(
		MaxBrokerCount,
		RestaterCheckInterval,
		BrokerReadyCheckInterval,
		BaseDir,
		SeedBatchFile,
		SimulatorFile,
		logger,
	)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Server",
		})
	})

	r.GET("/get_server", service.GetServer)
	r.POST("/init_broker", service.InitBrokerStart)
	r.POST("/done", service.SetBrokerWorkDone)
	r.POST("/start", service.Start)
	r.POST("/simulator/done", service.SetSimulatorDone)
	r.GET("/get_broker_server_ready", service.GetBrokerServerReady)

	r.LoadHTMLGlob("templates/*")
	r.Run() // default port :8080
}
