package saturn

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/drpaneas/ps5amazon/pkg/util"
)

var url string = "https://www.saturn.de/de/product/_sony-playstation%C2%AE5-digital-edition-2661939.html?utm_source=easymarketing&utm_medium=aff-content&utm_term=50201-912440140151390720&utm_campaign=AFF_Startseite&emid=5fbaca9a39f507795379c531"
var gatewayTimeout bool = false

func getHTML(page string) (doc *goquery.Document) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Request the HTML page.
	res, err := client.Get(page)
	if err != nil {
		gatewayTimeout = true
		return doc
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
		gatewayTimeout = true
		return doc
	}
	return doc
}

func isAvailable(s string) bool {
	if strings.Contains(s, "Dieser Artikel ist aktuell nicht verfügbar.") {
		return false
	}
	return true
}

func getPS5DigitalEdition(doc *goquery.Document) (text string, err error) {
	productQuery := fmt.Sprintf("#root > div.indexstyled__StyledAppWrapper-sc-1hu9cx8-0.cPLDhA > div.ProductDetailPagestyled__StyledPdpWrapper-sc-5s3nfq-1.hjoxyt > div:nth-child(1) > div > div.Cellstyled__StyledCell-sc-1wk5bje-0.ibdyBk.ProductDetailPagestyled__StyledPdpDetailCell-sc-5s3nfq-4.gLozy > div > div > div.Row__StyledRow-x4c83j-0.eaomqX.ProductDetailsstyled__StyledProductDetailRow-sc-12m2uf1-0.cJgftB > div > div > div > div > h4")
	doc.Find(productQuery).Each(func(i int, s *goquery.Selection) {
		text = util.ApplyTextFormat(s.Text())
	})
	if text == "" {
		err = fmt.Errorf("couldn't parse shit")
	}
	return text, err
}

// IsReadyToBuy returns true when PS5 is back in stock at saturn.de
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
	log.Println("[Saturn]     Sony PlayStation 5 - Digital Edition: " + ps5DigitalEdition)
	return false
}
