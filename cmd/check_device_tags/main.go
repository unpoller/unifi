package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/unpoller/unifi/v5"
)

func main() {
	apiKey := os.Getenv("UNIFI_API_KEY")
	if apiKey == "" {
		fmt.Fprintf(os.Stderr, "Error: UNIFI_API_KEY not set\n")
		os.Exit(1)
	}

	controllerURL := os.Getenv("UNIFI_URL")
	if controllerURL == "" {
		controllerURL = "https://192.168.1.1"
	}

	site := os.Getenv("UNIFI_SITE")
	if site == "" {
		site = "default"
	}

	config := &unifi.Config{
		APIKey:    apiKey,
		URL:       controllerURL,
		ErrorLog:  func(_ string, _ ...interface{}) {},
		DebugLog:  func(_ string, _ ...interface{}) {},
		Timeout:   30 * time.Second,
		VerifySSL: false,
	}

	uni, err := unifi.NewUnifi(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing UniFi client: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		_ = uni.Logout()
	}()

	// Fetch device tags
	path := fmt.Sprintf("/proxy/network/v2/api/site/%s/device-tags", site)

	body, err := uni.GetJSON(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching device tags: %v\n", err)
		os.Exit(1)
	}

	var tags []struct {
		ID               string   `json:"_id"`
		Name             string   `json:"name"`
		MemberDeviceMacs []string `json:"member_device_macs"`
	}

	if err := json.Unmarshal(body, &tags); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing device tags: %v\n", err)
		os.Exit(1)
	}

	// Find "test" tag
	fmt.Println("=== Device Tags ===")

	testTagFound := false

	for _, tag := range tags {
		if tag.Name == "test" {
			testTagFound = true

			fmt.Printf("\nTag: %s (ID: %s)\n", tag.Name, tag.ID)
			fmt.Printf("Device MACs with this tag: %d devices\n", len(tag.MemberDeviceMacs))

			for i, mac := range tag.MemberDeviceMacs {
				fmt.Printf("  %d. %s\n", i+1, mac)
			}
		}
	}

	if !testTagFound {
		fmt.Println("\nNo 'test' tag found.")
		fmt.Println("\nAvailable tags:")

		for _, tag := range tags {
			fmt.Printf("  - %s (%d devices)\n", tag.Name, len(tag.MemberDeviceMacs))
		}
	}

	// Also fetch devices to show names
	fmt.Println("\n=== Fetching device names ===")

	sites, err := uni.GetSites()
	if err != nil {
		fmt.Printf("Warning: Could not fetch sites: %v\n", err)

		return
	}

	var siteObj *unifi.Site

	for _, s := range sites {
		if s.Name == site {
			siteObj = s

			break
		}
	}

	if siteObj == nil {
		fmt.Printf("Warning: Site '%s' not found\n", site)

		return
	}

	devices, err := uni.GetDevices([]*unifi.Site{siteObj})
	if err != nil {
		fmt.Printf("Warning: Could not fetch devices: %v\n", err)

		return
	}

	// Create MAC to device name map
	macToName := make(map[string]string)
	for _, device := range devices.UAPs {
		macToName[device.Mac] = device.Name
	}

	for _, device := range devices.USWs {
		macToName[device.Mac] = device.Name
	}

	for _, device := range devices.UDMs {
		macToName[device.Mac] = device.Name
	}

	for _, device := range devices.USGs {
		macToName[device.Mac] = device.Name
	}

	for _, device := range devices.UXGs {
		macToName[device.Mac] = device.Name
	}

	for _, device := range devices.PDUs {
		macToName[device.Mac] = device.Name
	}

	// Show all devices with their tags
	fmt.Println("\n=== All Devices by Tag ===")

	for _, tag := range tags {
		fmt.Printf("\nTag: %s (%d devices)\n", tag.Name, len(tag.MemberDeviceMacs))

		for _, mac := range tag.MemberDeviceMacs {
			name := macToName[mac]
			if name == "" {
				name = "Unknown Device"
			}
			// Try to find device type
			deviceType := "Unknown"

			for _, d := range devices.UAPs {
				if d.Mac == mac {
					deviceType = "UAP"

					break
				}
			}

			if deviceType == "Unknown" {
				for _, d := range devices.USWs {
					if d.Mac == mac {
						deviceType = "USW"

						break
					}
				}
			}

			if deviceType == "Unknown" {
				for _, d := range devices.UDMs {
					if d.Mac == mac {
						deviceType = "UDM"

						break
					}
				}
			}

			if deviceType == "Unknown" {
				for _, d := range devices.USGs {
					if d.Mac == mac {
						deviceType = "USG"

						break
					}
				}
			}

			if deviceType == "Unknown" {
				for _, d := range devices.UXGs {
					if d.Mac == mac {
						deviceType = "UXG"

						break
					}
				}
			}

			if deviceType == "Unknown" {
				for _, d := range devices.PDUs {
					if d.Mac == mac {
						deviceType = "PDU"

						break
					}
				}
			}

			fmt.Printf("  - %s [%s] (%s)\n", name, deviceType, mac)
		}
	}
}
