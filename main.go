package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bendahl/uinput"
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

var stringToKeyCode = map[string]int{
	"ctrl":  uinput.KeyLeftctrl,
	"alt":   uinput.KeyLeftalt,
	"shift": uinput.KeyLeftshift,
	"super": uinput.KeyLeftmeta,
	"c":     uinput.KeyC,
	"v":     uinput.KeyV,
	"t":     uinput.KeyT,
	"f4":    uinput.KeyF4,
	" ":     uinput.KeySpace,
	"up":    uinput.KeyUp,
	"down":  uinput.KeyDown,
}

func pressKeys(keyboard uinput.Keyboard, action string) error {
	parts := strings.Split(strings.ToLower(action), "+")
	keyCodes := make([]int, 0, len(parts))

	for _, part := range parts {
		if code, ok := stringToKeyCode[part]; ok {
			keyCodes = append(keyCodes, code)
		} else {
			return fmt.Errorf("unknown key in action: '%s'", part)
		}
	}

	if len(keyCodes) == 0 {
		return fmt.Errorf("no valid keycodes found in action: '%s'", action)
	}

	// Press modifier keys (all keys except the last one)
	for i := 0; i < len(keyCodes)-1; i++ {
		err := keyboard.KeyDown(keyCodes[i])
		if err != nil {
			return fmt.Errorf("failed to press down key: %v", err)
		}
	}

	// Press and release the main key
	mainKey := keyCodes[len(keyCodes)-1]
	err := keyboard.KeyPress(mainKey)
	if err != nil {
		return fmt.Errorf("failed to press key: %v", err)
	}

	// Release modifier keys in reverse order
	for i := len(keyCodes) - 2; i >= 0; i-- {
		err := keyboard.KeyUp(keyCodes[i])
		if err != nil {
			return fmt.Errorf("failed to release key: %v", err)
		}
	}

	return nil
}

func main() {
	// Load configuration
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Find the ShuttleXpress device
	devicePath, err := findDevice()
	if err != nil {
		log.Fatalf("Could not find ShuttleXpress device: %v", err)
	}

	// Open the input device
	device, err := evdev.Open(devicePath)
	if err != nil {
		log.Fatalf("Error opening device: %v", err)
	}

	log.Printf("Using device: %s", devicePath)

	// Create virtual keyboard
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("shuttlexpress-virtual-keyboard"))
	if err != nil {
		log.Println("Error creating virtual keyboard. Please ensure:")
		log.Println("1. The 'uinput' kernel module is loaded (`sudo modprobe uinput`).")
		log.Println("2. You have write permissions to /dev/uinput (e.g., `sudo chmod 0666 /dev/uinput` or add a udev rule).")
		log.Fatalf("Error: %v", err)
	}
	defer keyboard.Close()

	log.Println("ShuttleXpress driver started. Press Ctrl+C to exit.")

	var lastJog int32 = -1 // Initialize with an invalid value

	// Event loop
	for {
		event, err := device.ReadOne()
		if err != nil {
			log.Fatalf("Error reading event: %v", err)
		}

		handleEvent(event, config, &lastJog, keyboard)
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
			return nil
		}
		defer device.File.Close()

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

func handleEvent(event *evdev.InputEvent, config *Config, lastJog *int32, keyboard uinput.Keyboard) {
	switch event.Type {
	case evdev.EV_KEY:
		if event.Value == 1 { // Key press
			keyCode := fmt.Sprintf("%d", event.Code)
			if action, ok := config.Buttons[keyCode]; ok {
				log.Printf("Button %s pressed, action: %s", keyCode, action)
				err := pressKeys(keyboard, action)
				if err != nil {
					log.Printf("Error simulating key press for action '%s': %v", action, err)
				}
			}
		}
	case evdev.EV_REL:
		switch event.Code {
		case evdev.REL_DIAL: // Jog
			currentJog := event.Value
			if *lastJog != -1 {
				delta := currentJog - *lastJog
				var actionKey string

				if (delta > 0 && delta < 128) || (delta < -128) { // Right turn
					actionKey = "1"
				} else if (delta < 0 && delta > -128) || (delta > 128) { // Left turn
					actionKey = "-1"
				}

				if actionKey != "" {
					if action, ok := config.Jog[actionKey]; ok {
						log.Printf("Jog event, action: %s", action)
						err := pressKeys(keyboard, action)
						if err != nil {
							log.Printf("Error simulating key press for action '%s': %v", action, err)
						}
					}
				}
			}
			*lastJog = currentJog
		case evdev.REL_WHEEL: // Ring
			value := fmt.Sprintf("%d", event.Value)
			if action, ok := config.Ring[value]; ok {
				log.Printf("Ring event, action: %s. (Note: Ring actions are not yet mapped to key presses)", action)
				// err := pressKeys(keyboard, action)
				// if err != nil {
				// 	log.Printf("Error simulating key press for action '%s': %v", action, err)
				// }
			}
		}
	}
}