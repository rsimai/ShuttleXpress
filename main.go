package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gvalkov/golang-evdev"
)

const (
	vendorID  = 0x0b33
	productID = 0x0020
)

type Config struct {
	Buttons map[string]string `json:"buttons"`
	Jog     map[string]string `json:"jog"`
	Ring    map[string]string `json:"ring"`
}

func main() {
	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Find the ShuttleXpress device
	devicePath, err := findDevice()
	if err != nil {
		fmt.Printf("Could not find ShuttleXpress device: %v\n", err)
		os.Exit(1)
	}

	// Open the input device
	device, err := evdev.Open(devicePath)
	if err != nil {
		fmt.Printf("Error opening device: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using device: %s\n", devicePath)
	fmt.Println("ShuttleXpress driver started. Press Ctrl+C to exit.")

	var lastJog int32 = -1 // Initialize with an invalid value

	// Event loop
	for {
		event, err := device.ReadOne()
		if err != nil {
			fmt.Printf("Error reading event: %v\n", err)
			os.Exit(1)
		}

		handleEvent(event, config, &lastJog)
	}
}

func findDevice() (string, error) {
	var devicePath string

	err := filepath.Walk("/dev/input", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasPrefix(info.Name(), "event") {
			return nil
		}

		device, err := evdev.Open(path)
		if err != nil {
			// Ignore errors from devices we can't open
			return nil
		}

		if device.Vendor == vendorID && device.Product == productID {
			devicePath = path
			return filepath.SkipDir // Stop searching
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if devicePath == "" {
		return "", fmt.Errorf("ShuttleXpress not found")
	}

	return devicePath, nil
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func handleEvent(event *evdev.InputEvent, config *Config, lastJog *int32) {
	switch event.Type {
	case evdev.EV_KEY:
		if event.Value == 1 { // Key press
			keyCode := fmt.Sprintf("%d", event.Code)
			if action, ok := config.Buttons[keyCode]; ok {
				fmt.Printf("Button %s pressed, action: %s\n", keyCode, action)
			}
		}
	case evdev.EV_REL:
		switch event.Code {
		case evdev.REL_DIAL: // Jog
			currentJog := event.Value
			if *lastJog != -1 {
				// The jog wheel is an 8-bit counter (0-255) that wraps around.
				delta := currentJog - *lastJog
				var actionKey string

				if (delta > 0 && delta < 128) || (delta < -128) { // Right turn
					actionKey = "1"
				} else if (delta < 0 && delta > -128) || (delta > 128) { // Left turn
					actionKey = "-1"
				}

				if action, ok := config.Jog[actionKey]; ok {
					fmt.Printf("Jog event, action: %s\n", action)
				}
			}
			*lastJog = currentJog
		case evdev.REL_WHEEL: // Ring
			value := fmt.Sprintf("%d", event.Value)
			if action, ok := config.Ring[value]; ok {
				fmt.Printf("Ring event, action: %s\n", action)
			}
		}
	}
}
