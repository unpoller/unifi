package mocks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"time"

	"github.com/unpoller/unifi/v5"
)

type MockHTTPTestServer struct {
	Server *httptest.Server
	mocked *MockUnifi
}

func NewMockHTTPTestServer() *MockHTTPTestServer {
	mocked := NewMockUnifi()
	m := &MockHTTPTestServer{mocked: mocked}
	s := httptest.NewServer(m)
	m.Server = s

	return m
}

func convertPathToRegexPattern(s string) string {
	tmp := strings.ReplaceAll(strings.ReplaceAll(s, "%s", "[^/]+"), "%d", "[0-9]+")

	return fmt.Sprintf("(%s)?%s", unifi.APIPrefixNew, tmp)
}

// compile regexp matches to paths
var (
	apiRogueAP         = regexp.MustCompile(convertPathToRegexPattern(unifi.APIRogueAP))
	apiStatusPath      = regexp.MustCompile(convertPathToRegexPattern(unifi.APIStatusPath))
	apiEventPath       = regexp.MustCompile(convertPathToRegexPattern(unifi.APIEventPath))
	apiSiteList        = regexp.MustCompile(convertPathToRegexPattern(unifi.APISiteList))
	apiSiteDPI         = regexp.MustCompile(convertPathToRegexPattern(unifi.APISiteDPI))
	apiClientDPI       = regexp.MustCompile(convertPathToRegexPattern(unifi.APIClientDPI))
	apiClientPath      = regexp.MustCompile(convertPathToRegexPattern(unifi.APIClientPath))
	apiAllUserPath     = regexp.MustCompile(convertPathToRegexPattern(unifi.APIAllUserPath))
	apiNetworkPath     = regexp.MustCompile(convertPathToRegexPattern(unifi.APINetworkPath))
	apiDevicePath      = regexp.MustCompile(convertPathToRegexPattern(unifi.APIDevicePath))
	apiLoginPath       = regexp.MustCompile(convertPathToRegexPattern(unifi.APILoginPath))
	apiLoginPathNew    = regexp.MustCompile(convertPathToRegexPattern(unifi.APILoginPathNew))
	apiLogoutPath      = regexp.MustCompile(convertPathToRegexPattern(unifi.APILogoutPath))
	apiEventPathIDS    = regexp.MustCompile(convertPathToRegexPattern(unifi.APIEventPathIDS)) //nolint:revive
	apiEventPathAlarms = regexp.MustCompile(convertPathToRegexPattern(unifi.APIEventPathAlarms))
	apiAnomaliesPath   = regexp.MustCompile(convertPathToRegexPattern(unifi.APIAnomaliesPath))
	apiCommandPath     = regexp.MustCompile(convertPathToRegexPattern(unifi.APICommandPath))
	apiDevMgrPath      = regexp.MustCompile(convertPathToRegexPattern(unifi.APIDevMgrPath))
)

type errorResponse struct {
	Error string `json:"error"`
}

type dataWrapper struct {
	Data any `json:"data"`
}

func respondResultOrErr(w http.ResponseWriter, v any, err error, wrapWithDataAttribute bool) {
	if err != nil {
		log.Printf("[ERROR] Answering mock response err=%+v value=%+v\n", err, v)
	} else {
		log.Printf("[DEBUG] Answering mock response value=%+v\n", v)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Via", "unifi-mock-server")

	if err != nil {
		e := errorResponse{
			Error: err.Error(),
		}
		b, _ := json.Marshal(e)

		w.WriteHeader(500)

		_, _ = w.Write(b)

		return
	}

	if wrapWithDataAttribute {
		response := dataWrapper{Data: v}
		b, err := json.Marshal(response)

		if err != nil {
			e := errorResponse{
				Error: err.Error(),
			}
			b, _ := json.Marshal(e)

			w.WriteHeader(500)
			_, _ = w.Write(b)

			return
		}

		w.WriteHeader(200)

		_, _ = w.Write(b)

		return
	}

	// no data wrapper
	b, err := json.Marshal(v)
	if err != nil {
		e := errorResponse{
			Error: err.Error(),
		}
		b, _ := json.Marshal(e)

		w.WriteHeader(500)

		_, _ = w.Write(b)

		return
	}

	w.WriteHeader(200)
	_, _ = w.Write(b)
}

func (m *MockHTTPTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.URL.Path)
	log.Printf("[DEBUG] Received mock request path=%s\n", p)

	switch {
	case apiRogueAP.MatchString(p):
		aps, err := m.mocked.GetRogueAPs(nil)
		respondResultOrErr(w, aps, err, true)

		return
	case apiStatusPath.MatchString(p):
		s, err := m.mocked.GetServerData()
		respondResultOrErr(w, s, err, true)

		return
	case apiEventPath.MatchString(p):
		events, err := m.mocked.GetEvents(nil, time.Hour)
		respondResultOrErr(w, events, err, true)

		return
	case apiSiteList.MatchString(p):
		sites, err := m.mocked.GetSites()
		respondResultOrErr(w, sites, err, true)

		return
	case apiSiteDPI.MatchString(p):
		dpi, err := m.mocked.GetSiteDPI(nil)
		respondResultOrErr(w, dpi, err, true)

		return
	case apiClientDPI.MatchString(p):
		dpi, err := m.mocked.GetClientsDPI(nil)
		respondResultOrErr(w, dpi, err, true)

		return
	case apiClientPath.MatchString(p):
		clients, err := m.mocked.GetClients(nil)
		respondResultOrErr(w, clients, err, true)

		return
	case apiAllUserPath.MatchString(p):
		users, err := m.mocked.GetUsers(nil, 1)
		respondResultOrErr(w, users, err, true)

		return
	case apiNetworkPath.MatchString(p):
		networks, err := m.mocked.GetNetworks(nil)
		respondResultOrErr(w, networks, err, true)

		return
	case apiDevicePath.MatchString(p):
		device, err := m.mocked.GetDevices(nil)
		// we need to change the format response for devices.
		// it is an array of mixed types in a singular {"data": [...all]}
		devices := make([]any, 0)

		for _, d := range device.UAPs {
			devices = append(devices, d)
		}

		for _, d := range device.UDMs {
			devices = append(devices, d)
		}

		for _, d := range device.USGs {
			devices = append(devices, d)
		}

		for _, d := range device.USWs {
			devices = append(devices, d)
		}

		for _, d := range device.PDUs {
			devices = append(devices, d)
		}

		for _, d := range device.UXGs {
			devices = append(devices, d)
		}

		for _, d := range device.UBBs {
			devices = append(devices, d)
		}

		for _, d := range device.UCIs {
			devices = append(devices, d)
		}

		respondResultOrErr(w, devices, err, true)

		return
	case apiLoginPath.MatchString(p):
		err := m.mocked.Login()
		respondResultOrErr(w, nil, err, true)

		return
	case apiLoginPathNew.MatchString(p):
		err := m.mocked.Login()
		respondResultOrErr(w, nil, err, true)

		return
	case apiLogoutPath.MatchString(p):
		err := m.mocked.Logout()
		respondResultOrErr(w, nil, err, true)

		return
	case apiEventPathIDS.MatchString(p):
		ids, err := m.mocked.GetIDS(nil, time.Now())
		respondResultOrErr(w, ids, err, true)

		return
	case apiEventPathAlarms.MatchString(p):
		alarms, err := m.mocked.GetAlarms(nil)
		respondResultOrErr(w, alarms, err, true)

		return
	case apiAnomaliesPath.MatchString(p):
		anomalies, err := m.mocked.GetAnomalies(nil, time.Now())
		respondResultOrErr(w, anomalies, err, true)

		return
	case apiDevMgrPath.MatchString(p):
		// todo
		w.WriteHeader(501)

		return
	case apiCommandPath.MatchString(p):
		// todo
		w.WriteHeader(501)

		return
	default:
		log.Println("[DEBUG] Answering mock response err=404 not found")
		http.NotFoundHandler().ServeHTTP(w, r)

		return
	}
}
