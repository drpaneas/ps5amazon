package main

import (
	"fmt"
	"time"

	"github.com/drpaneas/ps5amazon/pkg/amazon"
	"github.com/drpaneas/ps5amazon/pkg/mediamarkt"
	"github.com/drpaneas/ps5amazon/pkg/saturn"
	"github.com/drpaneas/ps5amazon/pkg/util"
)

func main() {
	// Loop until it's ready to buy
	for {
		if amazon.IsReadyToBuy() {
			util.Alarm()
			break
		}
		if mediamarkt.IsReadyToBuy() {
			util.Alarm()
			break
		}
		if saturn.IsReadyToBuy() {
			util.Alarm()
			break
		}
		// if gamestop.IsReadyToBuy() {
		// 	util.Alarm()
		// 	break
		// }
		// if alternate.IsReadyToBuy() {
		// 	util.Alarm()
		// 	break
		// }
		// if euronics.IsReadyToBuy() {
		// 	util.Alarm()
		// 	break
		// }
		fmt.Println("¯\\_(ツ)_/¯")
		time.Sleep(60 * time.Second)
	}
}
