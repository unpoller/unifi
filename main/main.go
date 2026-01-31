package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
)

func GetEnvString(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return strings.TrimSpace(value)
	}

	log.Printf("Environment variable %s not set, using fallback: %s\n", key, fallback)

	return fallback
}

func GetEnvInt64(key string, fallback int64) int64 {
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

func GetEnvInt(key string, fallback int) int {
	return int(GetEnvInt64(key, int64(fallback)))
}

func ShowResponse[T any](prefix string, data []T, numResponses int) {
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

// discoverConfig is the JSON shape for --config when using --discover.
// Use the same URL, user, and pass as your unpoller config.
type discoverConfig struct {
	URL    string `json:"url"`
	User   string `json:"user"`
	Pass   string `json:"pass"`
	APIKey string `json:"api_key"`
}

func loadConfig(path string) (*unifi.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	var c discoverConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	c.URL = strings.TrimSpace(strings.TrimSuffix(c.URL, "/"))
	if c.URL == "" || c.User == "" {
		return nil, fmt.Errorf("config must set url and user")
	}

	cfg := &unifi.Config{
		URL:      c.URL,
		User:     c.User,
		Pass:     c.Pass,
		APIKey:   c.APIKey,
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	return cfg, nil
}

func main() {
	discover := flag.Bool("discover", false, "discover API endpoints and write a shareable report (use with --config and --output)")
	configPath := flag.String("config", "", "path to JSON config for --discover (keys: url, user, pass); if unset, uses GOLIFT_UNIFI_* env")
	outputPath := flag.String("output", "api_endpoints_discovery.md", "path for discovery report when using --discover")

	flag.Parse()

	if *discover {
		var cfg *unifi.Config

		if *configPath != "" {
			var err error

			cfg, err = loadConfig(*configPath)
			if err != nil {
				log.Fatalln("FATAL:", err)
			}
		} else {
			cfg = &unifi.Config{
				URL:      GetEnvString("GOLIFT_UNIFI_URL", "http://localhost:8080"),
				User:     GetEnvString("GOLIFT_UNIFI_USER", "admin"),
				Pass:     GetEnvString("GOLIFT_UNIFI_PASS", ""),
				ErrorLog: log.Printf,
				DebugLog: log.Printf,
			}
		}

		uni, err := unifi.NewUnifi(cfg)
		if err != nil {
			log.Fatalln("FATAL:", err)
		}

		sites, err := uni.GetSites()
		if err != nil {
			log.Fatalln("FATAL:", err)
		}

		site := "default"
		if len(sites) > 0 && sites[0].Name != "" {
			site = sites[0].Name
		}

		if err := uni.DiscoverEndpoints(site, *outputPath); err != nil {
			log.Fatalln("FATAL:", err)
		}

		log.Printf("Discovery report written to %s (share this file with maintainers for API/404 issues).", *outputPath)

		return
	}

	// Default: use env and run normal data fetch
	var config = unifi.Config{
		User:     GetEnvString("GOLIFT_UNIFI_USER", "admin"),
		Pass:     GetEnvString("GOLIFT_UNIFI_PASS", ""),
		URL:      GetEnvString("GOLIFT_UNIFI_URL", "http://localhost:8080"),
		ErrorLog: log.Printf,
		DebugLog: log.Printf,
	}

	numResponses := GetEnvInt("GOLIFT_UNIFI_SHOW_RESPONSES", 0)

	uni, err := unifi.NewUnifi(&config)
	if err != nil {
		log.Fatalln(FATAL, err)
	}

	sites, err := uni.GetSites()
	if err != nil {
		log.Fatalln(FATAL, err)
	}

	ShowResponse("Sites", sites, numResponses)

	devices, err := uni.GetDevices(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("UAP Device", devices.UAPs, numResponses)
		ShowResponse("USG Device", devices.USGs, numResponses)
		ShowResponse("USW Device", devices.USWs, numResponses)
		ShowResponse("UDM Device", devices.UDMs, numResponses)
		ShowResponse("UXG Device", devices.UXGs, numResponses)
		ShowResponse("UCI Device", devices.UCIs, numResponses)
		ShowResponse("PDU Device", devices.PDUs, numResponses)
		ShowResponse("UBB Device", devices.UBBs, numResponses)
	}

	sitesDPI, err := uni.GetSiteDPI(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("SitesDPI", sitesDPI, numResponses)
	}

	clients, err := uni.GetClients(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("Clients", clients, numResponses)
	}

	clientsDPI, err := uni.GetClientsDPI(sites)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("ClientsDPI", clientsDPI, numResponses)
	}

	end := time.Now().UnixMilli()
	start := end - 3600000

	epochMillisTimePeriod := unifi.EpochMillisTimePeriod{
		StartEpochMillis: start,
		EndEpochMillis:   end,
	}

	clientTraffic, err := uni.GetClientTraffic(sites, &epochMillisTimePeriod, true)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("Client Traffic", clientTraffic, numResponses)
	}

	mac := GetEnvString("GOLIFT_UNIFI_MAC", "2c:d9:74:b8:13:46")

	clientTrafficByMac, err := uni.GetClientTrafficByMac(sites[0], &epochMillisTimePeriod, true, mac)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("Client Traffic By MAC", clientTrafficByMac, numResponses)
	}

	countryTraffic, err := uni.GetCountryTraffic(sites, &epochMillisTimePeriod)
	if err != nil {
		log.Println(ERROR, err)
	} else {
		ShowResponse("Country Traffic", countryTraffic, numResponses)
	}
}
