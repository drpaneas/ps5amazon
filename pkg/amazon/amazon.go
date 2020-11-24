package amazon

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/drpaneas/ps5amazon/pkg/util"
)

var amazonURL string = "https://www.amazon.de/s?k=ps5+digital+edition&i=videogames&crid=2OMYH102NXI6N&sprefix=ps5+digital+e%2Caps%2C166&ref=nb_sb_ss_c_2_13_ts-a-p"
var gatewayTimeout bool = false

// getAmazonHTML returns the HTML body for a given webpage
func getAmazonHTML(page string) (doc *goquery.Document) {

	client := &http.Client{Timeout: 5 * time.Second}

	// Request the HTML page.
	req, err := http.NewRequest("GET", amazonURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Inject the headers (stripped from Chrome > Developer Tools > Network tab) into the HTTP Request
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-language", "en-DE,en;q=0.9,el-GR;q=0.8,el;q=0.7,de-DE;q=0.6,de;q=0.5,en-US;q=0.4")
	req.Header.Set("cookie", util.AmazonCookie)
	req.Header.Set("downlink", "10")
	req.Header.Set("ect", "4g")
	req.Header.Set("rtt", "50")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"86\", \"\"Not\\A;Brand\";v=\"99\", \"Google Chrome\";v=\"86\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-use", "1")
	req.Header.Set("upgrade-insecure-requests", "50")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")

	// Send the HTTP Request and receive the HTTP Response
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		gatewayTimeout = true
		log.Printf("status code error: %d %s", res.StatusCode, res.Status)
		return doc
	}

	// Load the HTML <body>
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func isPS5DigitalEdition(s string) bool {
	if strings.Contains(s, "Sony PlayStation 5 - Digital Edition") {
		return true
	}
	return false
}

func getPS5DigitalEdition(doc *goquery.Document) (ps5 string, err error) {
	// Search for PS5 Digital Edition at amazon.de
	// The query has been taken from: Chrome Developer Tools > Right Click > Copy > Selector
	productQuery := fmt.Sprintf("#search > div.s-desktop-width-max.s-desktop-content.sg-row > div.sg-col-20-of-24.sg-col-28-of-32.sg-col-16-of-20.sg-col.sg-col-32-of-36.sg-col-8-of-12.sg-col-12-of-16.sg-col-24-of-28 > div > span:nth-child(4) > div.s-main-slot.s-result-list.s-search-results.sg-row > div")
	doc.Find(productQuery).Each(func(i int, s *goquery.Selection) {
		text := util.ApplyTextFormat(s.Text())
		if isPS5DigitalEdition(text) { // if you find "Sony PlayStation 5 - Digital Edition" that's the product we want
			ps5 = text
		}
	})
	if ps5 == "" {
		err = fmt.Errorf("Couldn't find PS5 Digital Edition")
	}
	return ps5, err
}

func getAmazonTitle(s string) string {
	tmp := strings.Split(s, "USK Rating")
	return tmp[0]
}

func getAmazonRating(s string) string {
	tmp := strings.SplitAfter(s, "USK Rating:")
	tmp2 := strings.Split(tmp[1], "|")
	return tmp2[0]
}

func getAmazonPrice(s string) string {
	tmp := strings.SplitAfter(s, "USK Rating:")
	tmp2 := strings.Split(tmp[1], "PlayStation 5")
	return tmp2[1]
}

func isAmazonAvailable(s string) bool {
	if strings.Contains(s, "Currently unavailable") {
		return false
	}
	return true
}

// IsReadyToBuy returns true when PS5 is back in stock at Amazon.de
func IsReadyToBuy() bool {
	// Fetch amazon page
	docHTML := getAmazonHTML(amazonURL)

	// Check if there is timeout
	if gatewayTimeout {
		gatewayTimeout = false
		return gatewayTimeout
	}

	// Check if Amazon scraping found PS5 Digital Edition
	ps5DigitalEdition, err := getPS5DigitalEdition(docHTML)
	if err != nil {
		log.Fatalf("Error: %v\nDescription: %s\n", err, ps5DigitalEdition)
	}

	// Check if PS5 is available to buy over there
	if isAmazonAvailable(ps5DigitalEdition) {
		fmt.Println("Buy now for " + getAmazonPrice(ps5DigitalEdition) + " !!!  --> " + amazonURL)
		return true
	}
	log.Println("[Amazon]     " + getAmazonTitle(ps5DigitalEdition) + getAmazonPrice(ps5DigitalEdition))
	return false
}
