package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server interface {
	GetServer(c *gin.Context)
	GetBrokerServerReady(c *gin.Context)
	InitBrokerStart(c *gin.Context)
	SetBrokerWorkDone(c *gin.Context)
	SetSimulatorDone(c *gin.Context)
	Start(c *gin.Context)
}

type Broker struct {
	Done bool `json:"done"`
}

type server struct {
	mutex           sync.Mutex        `json:"-"`
	ReRunTimes      int               `json:"rerun_times"`
	Brokers         map[string]Broker `json:"brokers"`
	Simulator       simulator         `json:"-"`
	SimulatorStatus string            `json:"simulator_status"`
	StartTime       time.Time         `json:"start_time"`
	Logger          *logrus.Logger    `json:"-"`

	MaxBrokerCount           int
	RestaterCheckInterval    time.Duration
	BrokerReadyCheckInterval time.Duration
	BaseDir                  string
	SeedBatchFile            string
	SimulatorFile            string
}

func NewServer(
	MaxBrokerCount int,
	RestaterCheckInterval time.Duration,
	BrokerReadyCheckInterval time.Duration,
	BaseDir string,
	SeedBatchFile string,
	SimulatorFile string,

	logger *logrus.Logger,
) Server {
	simulator := simulator{status: 0}
	_, simulatorStatus := simulator.GetStatus()
	brokers := make(map[string]Broker)
	return &server{
		mutex:           sync.Mutex{},
		Brokers:         brokers,
		ReRunTimes:      0,
		Simulator:       simulator,
		SimulatorStatus: simulatorStatus,
		Logger:          logger,

		MaxBrokerCount:           MaxBrokerCount,
		RestaterCheckInterval:    RestaterCheckInterval,
		BrokerReadyCheckInterval: BrokerReadyCheckInterval,
		BaseDir:                  BaseDir,
		SeedBatchFile:            SeedBatchFile,
		SimulatorFile:            SimulatorFile,
	}
}

func (s *server) GetServer(c *gin.Context) {
	_, simulatorStatus := s.Simulator.GetStatus()
	c.JSON(http.StatusOK, gin.H{
		"simulator_status": simulatorStatus,
		"rerun_times":      s.ReRunTimes,
		"broker_counts":    len(s.Brokers),
		"brokers":          s.Brokers,
		"uptime":           time.Since(s.StartTime).Seconds(),
	})
}

func (s *server) GetBrokerServerReady(c *gin.Context) {
	if len(s.Brokers) >= s.MaxBrokerCount {
		c.JSON(http.StatusOK, gin.H{"ready": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"ready": false})
	}
}

func (s *server) InitBrokerStart(c *gin.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	port := c.PostForm("port")
	s.Brokers[port] = Broker{Done: false}
	c.JSON(http.StatusOK, s.Brokers)
}

func (s *server) SetBrokerWorkDone(c *gin.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	port := c.PostForm("port")
	s.Brokers[port] = Broker{Done: true}
	c.JSON(http.StatusOK, s.Brokers)
}

func (s *server) SetSimulatorDone(c *gin.Context) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Simulator.SetDone()
	c.JSON(http.StatusOK, gin.H{"message": "done"})
}

func (s *server) Start(c *gin.Context) {
	s.start()
	go s.waitToRestart()
	go s.monitorLongRunning()
	c.JSON(http.StatusOK, gin.H{"message": "start"})
}

func (s *server) start() {
	// s.mutex.Lock()
	// defer s.mutex.Unlock()
	// run broker
	s.StartTime = time.Now()
	fmt.Println("Running broker CMD")
	err := s.runBrokerCMD()
	fmt.Println("Broker CMD done")
	if err != nil {
		s.Logger.WithFields(logrus.Fields{
			"message": "Failed to run broker",
			"error":   err.Error(),
		}).Error(s.ReRunTimes)
	} else {
		go s.waitBrokersReady()
	}
}

func (s *server) restart() {
	for i := 10; i > 0; i-- {
		fmt.Printf("Restarting in %d seconds\n", i)
		time.Sleep(1 * time.Second)
	}
	// s.mutex.Lock()
	s.Brokers = make(map[string]Broker)
	s.ReRunTimes++
	// s.mutex.Unlock()
	s.start()
}

func (s *server) waitToRestart() {
	// if all brokers are done, and simulator is done, then restart
	for {
		if len(s.Brokers) >= s.MaxBrokerCount {
			for _, broker := range s.Brokers {
				if broker.Done == false {
					break
				}
			}
			if s.Simulator.IsRunning() == false {
				s.restart()
			}
		}
		time.Sleep(s.RestaterCheckInterval)
	}
}

func (s *server) waitBrokersReady() {
	s.Simulator.SetWaiting()
	for {
		if len(s.Brokers) >= s.MaxBrokerCount {
			err := s.Simulator.Run(
				s.BaseDir,
				s.SimulatorFile,
			)
			if err != nil {
				s.Logger.WithFields(logrus.Fields{
					"message": "Failed to run simulator",
					"error":   err.Error(),
				}).Error(s.ReRunTimes)
			}
			break
		}
		fmt.Printf("Waiting for brokers to be ready, current count: %d\n", len(s.Brokers))
		time.Sleep(s.BrokerReadyCheckInterval)
	}
}

func (s *server) runBrokerCMD() error {
	cmdStr := fmt.Sprintf("cd %s && %s", s.BaseDir, s.SeedBatchFile)
	return RunCommand(cmdStr)
}

func (s *server) monitorLongRunning() {
	for {
		if time.Since(s.StartTime) > 3*time.Minute { // 30 minutes threshold for logging
			s.mutex.Lock()
			for port, broker := range s.Brokers {
				if !broker.Done {
					s.Logger.WithFields(logrus.Fields{
						"broker": port,
						"error":  "broker is not done",
					}).Error(
						s.ReRunTimes,
					)
				}
			}
			s.mutex.Unlock()
		}
		time.Sleep(5 * time.Minute) // Check every 5 minutes
	}
}
