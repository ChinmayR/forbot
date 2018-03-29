package forbot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Headers struct {
	contentType    string
	agent          string
	DatetimeFormat string
	auth           string
}

type OandaConnection struct {
	hostname       string
	port           int
	ssl            bool
	token          string
	accountID      string
	DatetimeFormat string
	headers        *Headers
}

const OANDA_AGENT string = "v20-golang/forbot-0.0.1"

func NewConnection(accountID string, token string, live bool) *OandaConnection {
	hostname := ""
	// should we use the live API?
	if live {
		hostname = "https://api-fxtrade.oanda.com/v3"
	} else {
		hostname = "https://api-fxpractice.oanda.com/v3"
	}

	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString("Bearer ")
	buffer.WriteString(token)

	authHeader := buffer.String()
	// Create headers for oanda to be used in requests
	headers := &Headers{
		contentType:    "application/json",
		agent:          OANDA_AGENT,
		DatetimeFormat: "RFC3339",
		auth:           authHeader,
	}
	// Create the connection object
	connection := &OandaConnection{
		hostname:  hostname,
		port:      443,
		ssl:       true,
		token:     token,
		headers:   headers,
		accountID: accountID,
	}

	return connection
}

func (c *OandaConnection) Request(endpoint string, params url.Values) []byte {
	client := http.Client{
		Timeout: time.Second * 5, // 5 sec timeout
	}

	generatedUrl := createUrl(c.hostname, endpoint, params)

	// New request object
	req, err := http.NewRequest(http.MethodGet, generatedUrl, nil)
	checkErr(err)

	body := makeRequest(c, endpoint, client, req)

	return body
}

func (c *OandaConnection) Send(endpoint string, data []byte, params url.Values) []byte {
	client := http.Client{
		Timeout: time.Second * 5, // 5 sec timeout
	}

	generatedUrl := createUrl(c.hostname, endpoint, params)

	// New request object
	req, err := http.NewRequest(http.MethodPost, generatedUrl, bytes.NewBuffer(data))
	checkErr(err)

	body := makeRequest(c, endpoint, client, req)

	return body
}

func (c *OandaConnection) Update(endpoint string, data []byte, params url.Values) []byte {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	generatedUrl := createUrl(c.hostname, endpoint, params)

	req, err := http.NewRequest(http.MethodPut, generatedUrl, bytes.NewBuffer(data))
	checkErr(err)
	body := makeRequest(c, endpoint, client, req)
	return body
}

func createUrl(host string, endpoint string, params url.Values) string {
	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(host)
	buffer.WriteString(endpoint)
	buffer.WriteString("?" + params.Encode())

	return buffer.String()
}

func makeRequest(c *OandaConnection, endpoint string, client http.Client, req *http.Request) []byte {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, getErr := client.Do(req)
	checkErr(getErr)
	body, readErr := ioutil.ReadAll(res.Body)
	checkErr(readErr)
	checkApiErr(body, endpoint)
	return body
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkApiErr(body []byte, route string) {
	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		log.SetFlags(log.LstdFlags | log.Llongfile)
		log.Fatal("\nOANDA API Error: " + bodyString + "\nOn route: " + route)
	}
}

func unmarshalJson(body []byte, data interface{}) {
	jsonErr := json.Unmarshal(body, &data)
	checkErr(jsonErr)
}
