package handlers

import (
	"fmt"
)

type simulator struct {
	status int
}

const (
	SimulatorWaiting = iota // 0
	SimulatorRunning        // 1
	SimulatorDone           // 2
)

func (s *simulator) GetStatus() (int, string) {
	var message string
	switch s.status {
	case 0:
		message = "Waiting for broker to be ready"
	case 1:
		message = "Simulator is running"
	case 2:
		message = "Simulator is done"
	}
	return s.status, message
}

func (s *simulator) IsRunning() bool {
	return s.status == SimulatorRunning
}

func (s *simulator) SetRunning() {
	s.status = 1
}

func (s *simulator) SetDone() {
	s.status = 2
}

func (s *simulator) SetWaiting() {
	s.status = 0
}

func (s *simulator) Run(
	BaseDir string,
	SimulatorFile string,
) error {
	cmdStr := fmt.Sprintf("cd %s && go run %s", BaseDir, SimulatorFile)
	err := RunCommand(cmdStr)
	if err != nil {
		s.SetRunning()
		return nil
	} else {
		return err
	}
}
