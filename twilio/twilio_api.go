package twilio

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	TWILIO_SID        = "ACd93e74ecb40ca3484c755f13683ac7a4"
	TWILIO_AUTH_TOKEN = "5d77325b645ff16f855dff823e693f93"
	TWILIO_URL_STR    = "https://api.twilio.com/2010-04-01/Accounts/" + TWILIO_SID + "/Messages.json"

	NUMBER_TO   = "+12063038987"
	NUMBER_FROM = "+12065390831"
)

func SendMsgFromData(data string) {
	msgData := url.Values{}
	msgData.Set("To", NUMBER_TO)
	msgData.Set("From", NUMBER_FROM)
	msgData.Set("Body", data)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", TWILIO_URL_STR, &msgDataReader)
	req.SetBasicAuth(TWILIO_SID, TWILIO_AUTH_TOKEN)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var respData map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&respData)
		if err == nil {
			fmt.Printf("Text sent: %s\n", respData["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}
