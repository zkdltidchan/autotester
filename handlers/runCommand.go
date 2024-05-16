package handlers

import (
	"fmt"
	"os/exec"
	"runtime"
)

func RunCommand(cmdStr string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", "start", "cmd.exe", "/c", cmdStr)
	} else if runtime.GOOS == "darwin" { // macOS
		cmd = exec.Command("osascript", "-e", `tell application "Terminal" to activate`, "-e", `tell application "Terminal" to do script "`+cmdStr+`"`)
	} else {
		return fmt.Errorf("unsupported platform")
	}

	fmt.Printf("Running command: %s\n", cmdStr)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %s\n", err)
		return err
	}
	return nil
}
