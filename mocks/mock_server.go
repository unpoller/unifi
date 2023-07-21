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

	"github.com/unpoller/unifi"
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
	tmp := strings.ReplaceAll(strings.ReplaceAll(s, "%s", "[a-zA-Z-_,.]+"), "%d", "[0-9]+")
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
	apiEventPathIDS    = regexp.MustCompile(convertPathToRegexPattern(unifi.APIEventPathIDS))
	apiEventPathAlarms = regexp.MustCompile(convertPathToRegexPattern(unifi.APIEventPathAlarms))
	apiAnomaliesPath   = regexp.MustCompile(convertPathToRegexPattern(unifi.APIAnomaliesPath))
	apiCommandPath     = regexp.MustCompile(convertPathToRegexPattern(unifi.APICommandPath))
	apiDevMgrPath      = regexp.MustCompile(convertPathToRegexPattern(unifi.APIDevMgrPath))
)

func respondResultOrErr(w http.ResponseWriter, v any, err error) {
	log.Printf("[DEBUG] Answering mock response err=%+v value=%+v\n", err, v)
	if err != nil {
		b, _ := json.Marshal(err)
		w.WriteHeader(500)
		_, _ = w.Write(b)
		return
	}
	b, _ := json.Marshal(v)
	w.WriteHeader(200)
	_, _ = w.Write(b)
}

func (m *MockHTTPTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.URL.Path)
	log.Printf("[DEBUG] Received mock request path=%s\n", p)
	switch {
	case apiRogueAP.MatchString(p):
		aps, err := m.mocked.GetRogueAPs(nil)
		respondResultOrErr(w, aps, err)
		return
	case apiStatusPath.MatchString(p):
		s, err := m.mocked.GetServerData()
		respondResultOrErr(w, s, err)
		return
	case apiEventPath.MatchString(p):
		events, err := m.mocked.GetEvents(nil, time.Hour)
		respondResultOrErr(w, events, err)
		return

	case apiSiteList.MatchString(p):
		sites, err := m.mocked.GetSites()
		respondResultOrErr(w, sites, err)
		return
	case apiSiteDPI.MatchString(p):
		dpi, err := m.mocked.GetSiteDPI(nil)
		respondResultOrErr(w, dpi, err)
		return
	case apiClientDPI.MatchString(p):
		dpi, err := m.mocked.GetClientsDPI(nil)
		respondResultOrErr(w, dpi, err)
		return
	case apiClientPath.MatchString(p):
		clients, err := m.mocked.GetClients(nil)
		respondResultOrErr(w, clients, err)
		return
	case apiAllUserPath.MatchString(p):
		users, err := m.mocked.GetUsers(nil, 1)
		respondResultOrErr(w, users, err)
		return
	case apiNetworkPath.MatchString(p):
		networks, err := m.mocked.GetNetworks(nil)
		respondResultOrErr(w, networks, err)
		return
	case apiDevicePath.MatchString(p):
		devices, err := m.mocked.GetDevices(nil)
		respondResultOrErr(w, devices, err)
		return
	case apiLoginPath.MatchString(p):
		err := m.mocked.Login()
		respondResultOrErr(w, nil, err)
		return
	case apiLoginPathNew.MatchString(p):
		err := m.mocked.Login()
		respondResultOrErr(w, nil, err)
		return
	case apiLogoutPath.MatchString(p):
		err := m.mocked.Logout()
		respondResultOrErr(w, nil, err)
		return
	case apiEventPathIDS.MatchString(p):
		ids, err := m.mocked.GetIDS(nil, time.Now())
		respondResultOrErr(w, ids, err)
		return
	case apiEventPathAlarms.MatchString(p):
		alarms, err := m.mocked.GetAlarms(nil)
		respondResultOrErr(w, alarms, err)
		return
	case apiAnomaliesPath.MatchString(p):
		anomalies, err := m.mocked.GetAnomalies(nil, time.Now())
		respondResultOrErr(w, anomalies, err)
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
