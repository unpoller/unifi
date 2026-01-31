// Package unifi: endpoint discovery for support and debugging.
// DiscoverEndpoints probes a set of known API paths and writes a shareable report.

package unifi

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// DiscoverResult holds the result of probing a single endpoint.
type DiscoverResult struct {
	Method string
	Path   string
	Status int
}

// DiscoverEndpoints probes known API paths on the controller and writes a markdown report
// to outputPath. Use the same config (URL, user, pass) as unpoller so users can share
// the output file for support. site is typically "default".
func (u *Unifi) DiscoverEndpoints(site, outputPath string) error {
	now := time.Now()
	end := now.UnixMilli()
	start := end - 3600000

	paths := []struct {
		method string
		path   string
	}{
		{"GET", APISiteList},
		{"GET", fmt.Sprintf(APIDevicePath, site)},
		{"GET", fmt.Sprintf(APIClientPath, site)},
		{"GET", fmt.Sprintf(APIEventPath, site)},
		{"GET", fmt.Sprintf(APIEventPathAlarms, site)},
		{"GET", fmt.Sprintf(APINetworkPath, site)},
		{"GET", fmt.Sprintf(APISiteDPI, site)},
		{"GET", fmt.Sprintf(APIClientDPI, site)},
		{"GET", fmt.Sprintf(APIAnomaliesPath, site)},
		{"GET", fmt.Sprintf(APISystemLogPath, site)},
		{"GET", fmt.Sprintf(APIDeviceTagsPath, site)},
		{"GET", fmt.Sprintf(APIActiveDHCPLeasesPath, site)},
		{"GET", fmt.Sprintf(APIWANEnrichedConfigPath, site)},
		{"GET", fmt.Sprintf(APIWANLoadBalancingStatusPath, site)},
		{"GET", fmt.Sprintf(APIWANLoadBalancingConfigPath, site)},
		{"GET", fmt.Sprintf(APIWANSLAsPath, site)},
		{"GET", fmt.Sprintf(APIClientTrafficPath, site, start, end, false)},
		{"GET", fmt.Sprintf(APICountryTrafficPath, site, start, end)},
		{"GET", fmt.Sprintf(APIAggregatedDashboard, site, 86400)},
		{"GET", fmt.Sprintf(APIRogueAP, site)},
		{"GET", fmt.Sprintf(APIAllUserPath, site)},
		{"GET", fmt.Sprintf(APIEventPathIDS, site)},
	}

	results := make([]DiscoverResult, 0, len(paths))

	for _, p := range paths {
		status, err := u.Probe(p.path)
		if err != nil {
			u.ErrorLog("discover: probe %s: %v", p.path, err)

			results = append(results, DiscoverResult{Method: p.method, Path: p.path, Status: -1})
		} else {
			results = append(results, DiscoverResult{Method: p.method, Path: p.path, Status: status})
		}
	}

	// WAN ISP status needs a wan id; skip or use a placeholder - try common "wan" id
	wanPaths := []string{"wan", "wan1", "wan2"}
	for _, w := range wanPaths {
		path := fmt.Sprintf(APIWANISPStatusPath, site, w)

		status, err := u.Probe(path)
		if err != nil {
			continue
		}

		results = append(results, DiscoverResult{Method: "GET", Path: path, Status: status})
	}

	return writeDiscoverReport(u.URL, site, results, outputPath)
}

func writeDiscoverReport(controllerURL, site string, results []DiscoverResult, outputPath string) error {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Path < results[j].Path
	})

	var b strings.Builder
	b.WriteString("# API Endpoints (discovery report)\n\n")
	b.WriteString(fmt.Sprintf("- **Controller**: %s\n", controllerURL))
	b.WriteString(fmt.Sprintf("- **Site**: %s\n", site))
	b.WriteString(fmt.Sprintf("- **Total endpoints probed**: %d\n", len(results)))
	b.WriteString("\n---\n\n")
	b.WriteString("| Method | Path | Status |\n")
	b.WriteString("|--------|------|--------|\n")

	for _, r := range results {
		statusStr := fmt.Sprintf("%d", r.Status)
		if r.Status < 0 {
			statusStr = "error"
		}

		pathEscaped := strings.ReplaceAll(r.Path, "|", "\\|")
		b.WriteString(fmt.Sprintf("| %s | `%s` | %s |\n", r.Method, pathEscaped, statusStr))
	}

	b.WriteString("\n---\n\n")
	b.WriteString("Share this file with maintainers when reporting API or 404 issues.\n")

	return os.WriteFile(outputPath, []byte(b.String()), 0o600)
}
