package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Time         map[string]string `json:"time"`
	Applications []string          `json:"applications"`
}

func main() {
	fmt.Println("Welcome to Stay Focused Daemon")
	configPath := filepath.Join(os.Getenv("HOME"), ".config", ".stay_focused.json")
	for {
		config := loadConfig(configPath)
		processesHandler(config)
		time.Sleep(3 * time.Second)
	}
}

func loadConfig(configPath string) *Config {
	config := &Config{}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(configPath), 0755)
		config.Applications = []string{}
		saveConfig(configPath, config)
	} else {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(data, config)
		if err != nil {
			panic(err)
		}
	}

	return config
}

func saveConfig(configPath string, config *Config) {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		panic(err)
	}
}

func parseTime(now time.Time, timeString string) time.Time {
	year, month, day := now.Date()
	parsedTime, _ := time.Parse("15:04", timeString)
	return time.Date(year, month, day, parsedTime.Hour(), parsedTime.Minute(), 0, 0, now.Location())
}

func processesHandler(config *Config) {
	now := time.Now()
	startTime := parseTime(now, config.Time["start"])
	endTime := parseTime(now, config.Time["end"])
	if now.After(startTime) && now.Before(endTime) {
		processes, err := listProcesses()
		if err != nil {
			fmt.Println("Error listing processes:", err)
			return
		}

		for _, process := range processes {
			for _, app := range config.Applications {
				matched, err := regexp.MatchString("(?i).*"+app+".*", process)
				if err != nil {
					fmt.Println("Error matching regex:", err)
					continue
				}
				if matched {
					killProcess(app)
				}
			}
		}
	}
}

func listProcesses() ([]string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("tasklist")
	case "darwin", "linux":
		cmd = exec.Command("ps", "-e")
	default:
		return nil, fmt.Errorf("unsupported platform")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	processes := strings.Split(string(output), "\n")
	return processes, nil
}

func killProcess(process string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/IM", process)
	case "darwin", "linux":
		cmd = exec.Command("pkill", "-f", process)
		showNotification("Stay Focused!", "Closing application: "+process)
	default:
		fmt.Println("unsupported platform")
		return
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error killing process:", err, " ---> ", process)
	} else {
		fmt.Println("Killed process:", process)
	}
}

func showNotification(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	cmd := exec.Command("osascript", "-e", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
