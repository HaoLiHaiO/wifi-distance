package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func getSignalStrength() (int, error) {
	var cmd *exec.Cmd

	// Detect OS and log output
	osType := detectOS()
	fmt.Println("Detected OS:", osType)

	switch osType {
	case "linux", "darwin":
		cmd = exec.Command("bash", "-c", "nmcli -t -f SIGNAL dev wifi | head -n 1")
	case "windows":
		cmd = exec.Command("powershell", "-Command", "(netsh wlan show interfaces) -match 'Signal' | ForEach-Object {$_ -replace '.*: ', ''}")
	default:
		return 0, fmt.Errorf("unsupported OS: %s", osType)
	}

	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	signalStr := strings.TrimSpace(string(out))
	fmt.Println("Raw Signal Strength Output:", signalStr)

	re := regexp.MustCompile(`\D`)
	cleanedSignalStr := re.ReplaceAllString(signalStr, "")

	signal, err := strconv.Atoi(cleanedSignalStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse signal strength: %v", err)
	}

	return signal, nil
}

func detectOS() string {
	rawOS := getOS()
	fmt.Println("Raw OS response:", rawOS)

	switch {
	case strings.Contains(strings.ToLower(rawOS), "windows"):
		return "windows"
	case strings.Contains(strings.ToLower(rawOS), "linux"):
		return "linux"
	case strings.Contains(strings.ToLower(rawOS), "darwin"):
		return "darwin"
	default:
		return "unknown"
	}
}

func getOS() string {
	if osName := os.Getenv("OS"); osName != "" {
		return osName
	}

	cmd := exec.Command("uname")
	out, err := cmd.Output()
	if err != nil {
		return "windows"
	}
	return strings.TrimSpace(string(out))
}

func estimateDistance(rssi int, frequency int) float64 {
	exp := (27.55 - (20 * math.Log10(float64(frequency))) + float64(rssi)) / 20.0
	return math.Pow(10, exp)
}

func main() {
	signal, err := getSignalStrength()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	frequency := 2400

	distance := estimateDistance(signal, frequency)
	fmt.Printf("Wi-Fi Signal Strength: %d dBm\n", signal)
	fmt.Printf("Estimated Distance from Router: %.2f meters\n", distance)
}
