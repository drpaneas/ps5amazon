package alternate

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/drpaneas/ps5amazon/pkg/util"
)

var url string = "https://www.alternate.de/Sony-Interactive-Entertainment/PlayStation-5-Digital-Edition-Spielkonsole/html/product/1651221"
var gatewayTimeout bool = false

func getHTML(page string) (doc *goquery.Document) {

	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("Accept-Language", "en-DE,en;q=0.9,el-GR;q=0.8,el;q=0.7,de-DE;q=0.6,de;q=0.5,en-US;q=0.4")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cookie", util.AlternateCookie)
	req.Header.Add("Host", "www.alternate.de")
	req.Header.Add("Referer", "https://www.alternate.de/PlayStation")
	req.Close = true

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

	// Load the HTML document
	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func isAvailable(s string) bool {
	if strings.Contains(s, "Artikel kann nicht gekauft werden") {
		return false
	}
	return true
}

func getPS5DigitalEdition(doc *goquery.Document) (text string, err error) {
	productQuery := fmt.Sprintf("#buyProduct > div.productMainContainerTable > div.productMainContainerRow > div.productShort > div.availability > p")
	doc.Find(productQuery).Each(func(i int, s *goquery.Selection) {
		text = util.ApplyTextFormat(s.Text())
	})
	if text == "" {
		err = fmt.Errorf("couldn't parse shit")
	}
	return text, err
}

// IsReadyToBuy returns true when PS5 is back in stock at Alternate.de
func IsReadyToBuy() bool {
	docHTML := getHTML(url)

	// Check if there is timeout
	if gatewayTimeout {
		gatewayTimeout = false
		return gatewayTimeout
	}

	// Check if scapring was ok
	ps5DigitalEdition, err := getPS5DigitalEdition(docHTML)
	if err != nil {
		log.Fatalf("Error: %v\nDescription: %s\n", err, ps5DigitalEdition)
	}

	// Check if PS5 is available to buy over there
	if isAvailable(ps5DigitalEdition) {
		fmt.Println("Buy now !!!  --> " + url)
		return true
	}
	log.Println("[Alternate]  Sony PlayStation 5 - Digital Edition: " + ps5DigitalEdition)
	return false
}
