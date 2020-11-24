package gamestop

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/drpaneas/ps5amazon/pkg/util"
)

var url string = "https://www.gamestop.de/PS5/Games/60315/playstation-5-digital-edition?tduid=131801cef2f394c45570fef3581636a3&utm_medium=affiliate&utm_source=1992115&utm_campaign=TradeDoubler_DE"
var gatewayTimeout bool = false

func getHTML(page string) (doc *goquery.Document) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Request the HTML page.
	res, err := client.Get(page)
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
	if strings.Contains(s, "Nicht verfügbar") {
		return false
	}
	return true
}

func getPS5DigitalEdition(doc *goquery.Document) (text string, err error) {
	productQuery := fmt.Sprintf("#prodMain > div.mainInfo > div.addCartBar > div.prodRightBlock > div.buySection > div.bigBuyButtons.SPNOpenMap > span > a")
	doc.Find(productQuery).Each(func(i int, s *goquery.Selection) {
		text = util.ApplyTextFormat(s.Text())
	})
	if text == "" {
		err = fmt.Errorf("couldn't parse shit")
	}
	return text, err
}

// IsReadyToBuy returns true when PS5 is back in stock at gamestop.de
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
	log.Println("[GameStop]   Sony PlayStation 5 - Digital Edition: " + ps5DigitalEdition)
	return false
}
