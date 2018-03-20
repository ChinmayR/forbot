package main

import (
	"fmt"
	"github.com/ChinmayR/forbot"
)

const (
	ACCOUNT_ID = "509983"
	AUTH_TOKEN = "91e0ecb7a2d464feb06769ee342b14b0-d8b4f92f7590add749377505842e6a69"
)

func main() {
	oandaCon := forbot.NewConnection(ACCOUNT_ID, AUTH_TOKEN, true)
	history := oandaCon.GetCandles("EUR_USD")
	fmt.Println(history)
}