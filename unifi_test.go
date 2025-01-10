package unifi // nolint: testpackage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnifi(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u := "http://127.0.0.1:64431"
	c := &Config{
		User:      "user1",
		Pass:      "pass2",
		URL:       u,
		VerifySSL: false,
		DebugLog:  discardLogs,
	}
	authReq, err := NewUnifi(c)
	a.NotNil(err)
	a.EqualValues(u, authReq.URL)
	a.Contains(err.Error(), "connection refused", "an invalid destination should produce a connection error.")
}

func TestNewUnifiAPIKey(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	u := "http://127.0.0.1:64431"
	c := &Config{
		APIKey:    "fakekey",
		URL:       u,
		VerifySSL: false,
		DebugLog:  discardLogs,
	}
	authReq, err := NewUnifi(c)
	a.NotNil(err)
	a.EqualValues(u, authReq.URL)
	a.Contains(err.Error(), "connection refused", "an invalid destination should produce a connection error.")
}

func TestUniReq(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := "/test/path"
	u := "http://some.url:8443"

	// Test empty parameters.
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{URL: u, DebugLog: discardLogs}}
	r, err := authReq.UniReq(p, "")

	a.Nil(err, "newrequest must not produce an error")
	a.EqualValues(p, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(u, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("GET", r.Method, "without parameters the method must be GET")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Test with parameters
	k := "key1=value9&key2=value7"
	authReq = &Unifi{Client: &http.Client{}, Config: &Config{URL: "http://some.url:8443", DebugLog: discardLogs}}
	r, err = authReq.UniReq(p, k)
	a.Nil(err, "newrequest must not produce an error")

	a.EqualValues(p, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(u, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("POST", r.Method, "with parameters the method must be POST")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Check the parameters.
	d, err := io.ReadAll(r.Body)
	a.Nil(err, "problem reading request body, POST parameters may be malformed")
	a.EqualValues(k, string(d), "POST parameters improperly encoded")
}

func TestUniReqPut(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	p := "/test/path"
	u := "http://some.url:8443"

	// Test empty parameters.
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{URL: u, DebugLog: discardLogs}}
	_, err := authReq.UniReqPut(p, "")
	a.NotNil(err, "empty params must produce an error")

	// Test with parameters
	k := "key1=value9&key2=value7"
	authReq = &Unifi{Client: &http.Client{}, Config: &Config{URL: "http://some.url:8443", DebugLog: discardLogs}}
	r, err := authReq.UniReqPut(p, k)
	a.Nil(err, "newrequest must not produce an error")

	a.EqualValues(p, r.URL.Path,
		"the provided apiPath was not added to http request")
	a.EqualValues(u, r.URL.Scheme+"://"+r.URL.Host, "URL improperly encoded")
	a.EqualValues("PUT", r.Method, "with parameters the method must be POST")
	a.EqualValues("application/json", r.Header.Get("Accept"), "Accept header must be set to application/json")

	// Check the parameters.
	d, err := io.ReadAll(r.Body)
	a.Nil(err, "problem reading request body, PUT parameters may be malformed")
	a.EqualValues(k, string(d), "PUT parameters improperly encoded")
}

func TestUnifiIntegrationAPIKeyInjected(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") == "fakekey" {
			w.WriteHeader(http.StatusOK)

			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}))
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{APIKey: "fakekey", URL: srv.URL, DebugLog: discardLogs}}
	authResp, err := authReq.UniReqPost("/test", "")
	a.Nil(err, "newrequest must not produce an error")
	a.EqualValues("POST", authResp.Method, "with parameters the method must be POST")
}

func TestUnifiIntegrationUserPassInjected(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.EqualFold(r.URL.Path, "/api/login") {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("error reading body:%v\n", err)

			return
		}

		type userPass struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var up userPass

		err = json.Unmarshal(data, &up)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Printf("error decoding body: %s: %s\n", string(data), err)

			return
		}

		if strings.EqualFold(up.Username, "fakeuser") && strings.EqualFold(up.Password, "fakepass") {
			w.WriteHeader(http.StatusOK)
		}

		w.WriteHeader(http.StatusUnauthorized)
	}))
	authReq := &Unifi{Client: &http.Client{}, Config: &Config{User: "fakeuser", Pass: "fakepass", URL: srv.URL, DebugLog: discardLogs}}
	err := authReq.Login()
	a.Nil(err, "user/pass login must not produce an error")
}

/* NOT DONE: OPEN web server, check parameters posted, more. These tests are incomplete.
a.EqualValues(`{"username": "user1","password": "pass2"}`, string(post_params),
	"user/pass json parameters improperly encoded")
*/
