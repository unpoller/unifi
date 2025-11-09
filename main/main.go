package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
)

func getEnvString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return strings.TrimSpace(value)
	}
	log.Printf("Environment variable %s not set, using fallback: %s\n", key, fallback)
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	value, ok := os.LookupEnv(key)
	if ok {
		value = strings.TrimSpace(value)
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Printf("Environment variable %s not set, using fallback: %d\n", key, fallback)
		} else {
			return i
		}
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	return int(getEnvInt64(key, int64(fallback)))
}

func show[T any](prefix string, data []T, numResponses int) {
	log.Printf("%d Unifi %s found.\n", len(data), prefix)
	if numResponses <= 0 {
		return
	}
	for i, r := range data {
		if i >= numResponses {
			break
		}
		log.Printf("Response %d:\n", i)
		jsonData, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			log.Fatalf("Error marshaling to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	}
}

const ERROR = "ERROR:"
const FATAL = "FATAL:"

func main() {
	var config = unifi.Config{
		User:     getEnvString("GOLIFT_UNIFI_USER", "admin"),
		Pass:     getEnvString("GOLIFT_UNIFI_PASS", ""),
		URL:      getEnvString("GOLIFT_UNIFI_URL", "http://localhost:8080"),
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	num_responses := getEnvInt("GOLIFT_UNIFI_SHOW_RESPONSES", 0)

	uni, err := unifi.NewUnifi(&config)
	if err != nil {
		log.Fatalln(FATAL, err)
	}

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln(FATAL, err)
	} else {
		show("Sites", sites, num_responses)
	}

	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("UAP Device", devices.UAPs, num_responses)
		show("USG Device", devices.USGs, num_responses)
		show("USW Device", devices.USWs, num_responses)
		show("UDM Device", devices.UDMs, num_responses)
		show("UXG Device", devices.UXGs, num_responses)
		show("UCI Device", devices.UCIs, num_responses)
		show("PDU Device", devices.PDUs, num_responses)
		show("UBB Device", devices.UBBs, num_responses)
	}

	sitesDPI, err := uni.GetSiteDPI(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("SitesDPI", sitesDPI, num_responses)
	}

	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("Clients", clients, num_responses)
	}

	clientsDPI, err := uni.GetClientsDPI(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("ClientsDPI", clientsDPI, num_responses)
	}

	end := time.Now().UnixMilli()
	start := end - 3600000

	epochMillisTimePeriod := unifi.EpochMillisTimePeriod{
		StartEpochMillis: start,
		EndEpochMillis:   end,
	}

	client_traffic, err := uni.GetClientTraffic(sites, &epochMillisTimePeriod, true)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("Client Traffic", client_traffic, num_responses)
	}

	mac := getEnvString("GOLIFT_UNIFI_MAC", "2c:d9:74:b8:13:46")
	client_traffic_by_mac, err := uni.GetClientTrafficByMac(sites[0], &epochMillisTimePeriod, true, mac)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("Client Traffic By MAC", client_traffic_by_mac, num_responses)
	}

	country_traffic, err := uni.GetCountryTraffic(sites, &epochMillisTimePeriod)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		show("Country Traffic", country_traffic, num_responses)
	}
}
