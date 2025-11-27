package main

import (
	"fmt"
	"os"
	"os/exec"
)

func checkBwCommand() (bool, string) {
	path, err := exec.LookPath("bw")
	if err != nil {
		return false, ""
	}
	return true, path
}

func main() {
	installed, path := checkBwCommand()
	if !installed {
		fmt.Fprintf(os.Stderr, "[ERROR] ❌ bw command is not installed...\n")
		os.Exit(1)
	}
	fmt.Printf("[INFO] ✅ bw command is installed! (path: %s)\n", path)

}

