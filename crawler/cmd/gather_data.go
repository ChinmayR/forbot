package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/ChinmayR/forbot/constants"
	"github.com/ChinmayR/forbot/crawler"
)

const initialUrl = "https://www.daytradingforexlive.com/forums/xfa-blogs/sterling.106/?page="

func main() {
	// open the file
	pointsFile, err := os.OpenFile(constants.FILENAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	allPoints := make([]*crawler.CrawledPoints, 0)

	// if file does not exist then crawl and write to the file else read straight from file
	if _, err := os.Stat(constants.FILENAME); os.IsNotExist(err) {
		pagesToGo := 5

		for i := 1; i <= pagesToGo; i++ {
			log.Println("Crawling " + initialUrl + strconv.Itoa(i))
			resp, err := crawler.GetPageOfDTFL(initialUrl + strconv.Itoa(i))
			if err != nil {
				log.Println(err)
				return
			}

			crawledPoints, err := crawler.TokenizeMainPage(resp, 500)
			if err != nil {
				log.Println(err)
				return
			}
			allPoints = append(allPoints, crawledPoints...)
		}

		b, err := json.Marshal(allPoints)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(b))
		pointsFile.Write(b)

		pointsFile.Close()
	} else {
		b, err := ioutil.ReadAll(pointsFile)
		if err != nil {
			log.Println(err)
			return
		}
		json.Unmarshal(b, &allPoints)
		log.Println(allPoints)
	}

}
