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
		ErrorLog:  func(msg string, args ...interface{}) {},
		DebugLog:  func(msg string, args ...interface{}) {},
		Timeout:   30 * time.Second,
		VerifySSL: false,
	}

	uni, err := unifi.NewUnifi(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing UniFi client: %v\n", err)
		os.Exit(1)
	}
	defer uni.Logout()

	// Get sites
	sites, err := uni.GetSites()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching sites: %v\n", err)
		os.Exit(1)
	}

	var siteObj *unifi.Site
	for _, s := range sites {
		if s.Name == site {
			siteObj = s
			break
		}
	}

	if siteObj == nil {
		fmt.Fprintf(os.Stderr, "Site '%s' not found\n", site)
		os.Exit(1)
	}

	// Fetch devices (should be automatically enriched with tags)
	fmt.Println("=== Fetching Devices with Tags ===")
	devices, err := uni.GetDevices([]*unifi.Site{siteObj})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching devices: %v\n", err)
		os.Exit(1)
	}

	// Show devices with their tags
	fmt.Printf("\n=== UAPs ===\n")
	for _, device := range devices.UAPs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== USWs ===\n")
	for _, device := range devices.USWs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== UDMs ===\n")
	for _, device := range devices.UDMs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== USGs ===\n")
	for _, device := range devices.USGs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== UXGs ===\n")
	for _, device := range devices.UXGs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== PDUs ===\n")
	for _, device := range devices.PDUs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== UBBs ===\n")
	for _, device := range devices.UBBs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	fmt.Printf("\n=== UCIs ===\n")
	for _, device := range devices.UCIs {
		tagsJSON, _ := json.Marshal(device.Tags)
		fmt.Printf("  %s (%s) - Tags: %s\n", device.Name, device.Mac, string(tagsJSON))
	}

	// Also test GetDeviceTags directly
	fmt.Println("\n=== Testing GetDeviceTags Method ===")
	tags, err := uni.GetDeviceTags(siteObj)
	if err != nil {
		fmt.Printf("Error fetching tags: %v\n", err)
	} else {
		fmt.Printf("Found %d tags:\n", len(tags))
		for _, tag := range tags {
			fmt.Printf("  - %s (%d devices)\n", tag.Name, len(tag.MemberDeviceMacs))
		}
	}
}
