package crawler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ChinmayR/forbot/twilio"
	"golang.org/x/net/html"
)

type CrawledPoints struct {
	DataUrl   string
	Date      time.Time
	ManPoints []ManPointsForSymbol
}

type ManPointsForSymbol struct {
	Symbol string
	Points []float64
}

const urlToCall = "https://www.daytradingforexlive.com/forums/xfa-blog-home/"

// only send a text for the crawled data if it is 10:30PM PST
func SendText(manPoints []ManPointsForSymbol, dataUrl string, date time.Time) {
	PST, _ := time.LoadLocation("America/Los_Angeles")
	curTimeInPST := time.Now().In(PST)
	if curTimeInPST.Hour() == 22 && curTimeInPST.Minute() >= 30 && curTimeInPST.Minute() < 45 {
		//check if the new points for today have not been uploaded yet, if not then append to text message
		newPointsSubmitted := true
		fmt.Println(curTimeInPST)
		fmt.Println(date)
		if date.Year() == curTimeInPST.Year() && date.Month() == curTimeInPST.Month() && date.Day() == curTimeInPST.Day() {
			newPointsSubmitted = false
		}
		stringToText := "Crawled successfully for " + date.Format("Jan 2 2006")
		if !newPointsSubmitted {
			stringToText += "(yesterdays points)"
		}
		stringToText += "\n"

		for _, manPoint := range manPoints {
			stringToText += manPoint.Symbol + ": "
			for _, point := range manPoint.Points {
				stringToText += fmt.Sprintf("%v, ", point)
			}
			// remove the last comma added while printing the points
			stringToText = stringToText[:len(stringToText)-2] + "\n"
		}
		stringToText += "\n" + dataUrl
		twilio.SendMsgFromData(stringToText, "")
	} else {
		log.Println("Didn't send text after crawling DTFL because its not between 10:25 and 10:40PM")
	}
}

var GlobalManPointsData []*CrawledPoints

func GetTodayManipulationPoints() ([]*CrawledPoints, error) {
	if GlobalManPointsData != nil {
		log.Println("Found cached crawled data so just returning that.")
		return GlobalManPointsData, nil
	}

	resp, err := GetPageOfDTFL(urlToCall)
	if err != nil {
		return nil, err
	}

	crawledPoints, err := TokenizeMainPage(resp)
	if err != nil {
		twilio.SendMsgFromData(fmt.Sprintf("Error crawling: %v", err.Error()[0:100]), "")
		return nil, err
	}

	sort.Slice(crawledPoints, func(i, j int) bool {
		return crawledPoints[j].Date.Before(crawledPoints[i].Date)
	})

	for i, crawledPoint := range crawledPoints {
		log.Printf("Crawler run output %v: %v\n", i, crawledPoint)
	}
	SendText(crawledPoints[0].ManPoints, crawledPoints[0].DataUrl, crawledPoints[0].Date)
	GlobalManPointsData = crawledPoints

	return crawledPoints, nil
}

/*
<li id="entry-1470" class="entry">
			<a href="members/sterling.106/" class="avatar Av106s" data-avatarhtml="true"><img src="data/avatars/s/0/106.jpg?1511237315" width="48" height="48" alt="Sterling" /></a>
			<h3><a href="xfa-blog-entry/daily-market-preview-april-16th-2017.1470/">Daily Market Preview - April 16th 2017</a></h3>
			<div class="message">
				<blockquote>

				</blockquote>
			</div>
		</li>
*/
func TokenizeMainPage(resp *http.Response) ([]*CrawledPoints, error) {
	z := html.NewTokenizer(resp.Body)

	retVal := make([]*CrawledPoints, 0)
	var wg sync.WaitGroup
	pagesToVisit := 5

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			// End of the document, we're done
			break
		}

		switch {
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" && strings.HasPrefix(a.Val, "xfa-blog-entry/daily-market-preview") && pagesToVisit > 0 {

						pagesToVisit--
						wg.Add(1)
						dailyPreviewUrl := "https://www.daytradingforexlive.com/forums/" + a.Val

						go func(url string, passedRetVal *[]*CrawledPoints, wg2 *sync.WaitGroup) {
							defer wg2.Done()
							log.Println("Making Call: " + dailyPreviewUrl)
							respTemp, err := GetPageOfDTFL(dailyPreviewUrl)
							if err != nil {
								log.Println("Error while getting webpage: " + err.Error())
								return
							}
							manPoints, innerPageDate, err := TokenizeInnerPage(respTemp)
							if err != nil {
								log.Println("Error while parsing inner webpage: " + err.Error())
								return
							}

							if innerPageDate == nil {
								log.Printf("ERROR: Returned no inner date for link %v so ignoring manpoints %v\n", dailyPreviewUrl, manPoints)
								return
							}

							crawlPoint := &CrawledPoints{}
							crawlPoint.DataUrl = url
							crawlPoint.ManPoints = manPoints
							crawlPoint.Date = *innerPageDate
							*passedRetVal = append(*passedRetVal, crawlPoint)

							log.Printf("Added: %v\n", crawlPoint)
							return

						}(dailyPreviewUrl, &retVal, &wg)

						break
					}
				}
			}
			break
		}
	}
	wg.Wait()

	if len(retVal) == 0 {
		return nil, errors.New("no daily market preview URL found to crawl inner")
	} else {
		return retVal, nil
	}
}

/*
	<blockquote class="ugc baseHtml"><iframe width="500" height="300" src="https://www.youtube.com/embed/QKfyFG4jegc?wmode=opaque" frameborder="0" allowfullscreen></iframe><br />
<br />
<b>EUR/USD - </b>Looking for the second push down.<br />
<b>Manipulation Point/s -</b> 1.2350<br />
<br />
<b>EUR/JPY Manipulation Point/s - </b>132.85 Upper - 132.31 &amp; 131.85 &amp; 131.12 Lower<br />
<br />
-Sterling<br />
<br />
<b>MAX STOP:</b><br />
<b>EUR/USD:</b> 20 Pips<br />
<b>EUR/JPY:</b> 20 Pips</blockquote>

******For Date:
<h3 class="title customizeTitle">
			<div class="entryDetails">Views: 113</div>
			Daily Market Preview - June 8th 2018
			<span class="datetime muted"><abbr class="DateTime" data-time="1528405885" data-diff="253288" data-datestring="Jun 7, 2018" data-timestring="10:11 PM">Jun 7, 2018 at 10:11 PM</abbr></span>
</h3>

*/
func TokenizeInnerPage(resp *http.Response) ([]ManPointsForSymbol, *time.Time, error) {
	hm := make(map[string]ManPointsForSymbol)
	retManPoints := make([]ManPointsForSymbol, 0)

	z := html.NewTokenizer(resp.Body)

	var parsedDate *time.Time

L:
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			break L
		case tt == html.TextToken:
			t := z.Token()
			if strings.HasPrefix(t.Data, "Daily Market Preview") {
				dateTokens := strings.Split(strings.TrimSpace(strings.Split(t.Data, "-")[1]), " ")
				monthToken := dateTokens[0][:3]
				dayToken := dateTokens[1][0 : len(dateTokens[1])-2]
				dateString := monthToken + " " + dayToken + " " + dateTokens[2]
				log.Printf("Parsing date %v\n", dateString)

				date, err := time.Parse("Jan 2 2006", dateString)
				if err != nil {
					log.Printf("Error parsing date %v\n", dateString)
					break L
				}
				//date = date.Add(24 * time.Hour) // the daily market review is always posted for the next day
				log.Printf("Parsed date %v as %v\n", dateString, date.String())
				parsedDate = &date
			}
			break
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "blockquote" {
				currentSymbol := ""
			LL:
				for {
					ttt := z.Next()
					tempToken := z.Token()
					if strings.Contains(tempToken.Data, "MAX STOP") {
						break LL
					}
					switch {
					case ttt == html.ErrorToken:
						// End of the document, we're done
						break LL
					case ttt == html.TextToken:
						splitSting := strings.Split(tempToken.Data, " ")
						for _, str := range splitSting {
							str = strings.TrimSpace(str)
							if strings.Contains(str, "/") && strings.Index(str, "/") == 3 && len(str) == 7 {
								log.Println("Parsing symbol: " + str)
								currentSymbol = str
								hm[currentSymbol] = ManPointsForSymbol{Symbol: currentSymbol, Points: make([]float64, 0)}
							}
							floatVal, err := strconv.ParseFloat(str, 64)
							if currentSymbol != "" && err == nil {
								manPoints := hm[currentSymbol]
								manPoints.Points = append(manPoints.Points, floatVal)
								hm[currentSymbol] = manPoints
							}
						}
					}
				}

				for _, val := range hm {
					val.Symbol = strings.Replace(val.Symbol, "/", "_", -1)
					retManPoints = append(retManPoints, val)
				}
				return retManPoints, parsedDate, nil
			}
		}
	}

	return nil, nil, errors.New("no daily market preview URL found to crawl inner")
}

func Login(client *http.Client) {
	urlOfPage := "https://www.daytradingforexlive.com/forums/login/login"
	data := url.Values{}
	data.Set("login", "ChinmayR1992")
	data.Add("password", "daytradingforexlive")
	data.Add("cookie_check", "1")
	data.Add("_xfToken", "")
	data.Add("redirect", "https://www.daytradingforexlive.com/forums/")

	r, _ := http.NewRequest("POST", urlOfPage, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	//r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//r.Header.Add("Referer", "https://www.daytradingforexlive.com/forums/login/")
	//r.Header.Add("Upgrade-Insecure-Requests", "1")
	//r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	//r.Header.Add("Accept-Encoding", "gzip, deflate, br")

	GetFirstSession(client) //the cookie gets stored in the client jar
	resp, _ := client.Do(r)
	defer resp.Body.Close()
}

func GetPageOfDTFL(urlOfPage string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: jar}
	if err != nil {
		return nil, err
	}
	Login(client) // the cookie will get stored in the client jar

	req, err := http.NewRequest("GET", urlOfPage, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	//req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Add("Referer", "https://www.daytradingforexlive.com/forums/login/")
	//req.Header.Add("Upgrade-Insecure-Requests", "1")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Got status code %v for url %v", resp.StatusCode, urlOfPage))
	}
	return resp, nil
}

func GetFirstSession(client *http.Client) string {
	urlOfPage := "https://www.daytradingforexlive.com/forums/login"
	r, _ := http.NewRequest("GET", urlOfPage, nil) // URL-encoded payload
	resp, _ := client.Do(r)
	retVal := strings.Split(strings.Split(resp.Header.Get("Set-Cookie"), ";")[0], "=")[1]
	return retVal
}
