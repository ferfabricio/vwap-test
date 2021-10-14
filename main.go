package main

import (
	"fmt"
	"strconv"

	"github.com/ferfabricio/vwap-calculation-engine/domain/currencies"
	"github.com/ferfabricio/vwap-calculation-engine/domain/vwap"
	"github.com/ferfabricio/vwap-calculation-engine/infrastructure/coinbase"
)

func main() {
	c, err := coinbase.NewClient()
	if err != nil {
		panic(err)
	}
	// Add the product ids that must be subscribed on Coinbase
	err = c.Configure([]string{
		currencies.BitcoinDollar,
		currencies.EthereumBitcoin,
		currencies.EthereumDollar,
	})

	if err != nil {
		panic(err)
	}

	// Create a notification channel
	nf := make(chan vwap.CalculationEvent, 200)

	// Create a new VWAP calculation from domain
	v := vwap.NewWithNotification(nf)

	// Add all pairs to be calculated
	v.AddPair(currencies.BitcoinDollar)
	v.AddPair(currencies.EthereumBitcoin)
	v.AddPair(currencies.EthereumDollar)

	// Create the result channel to publish Coinbase trading pairs
	ch := make(chan coinbase.MatchMessage, 200)

	// Start collect trading pairs
	go c.GetData(ch)
	defer close(ch)
	for {
		select {
		case r := <-ch:
			// If the received message is of the type match add trade to VWAP calculator
			if r.Type == coinbase.MatchType {
				p, err := strconv.ParseFloat(r.Price, 32)
				if err != nil {
					panic(err)
				}
				s, err := strconv.ParseFloat(r.Size, 32)
				if err != nil {
					panic(err)
				}
				if err := v.AddTrade(r.ProductId, p, s); err != nil {
					panic(err)
				}
			}
		case n := <-nf:
			// Print notifications
			fmt.Printf("Pair: %s, Price: %f, Quantity: %f, VWAP: %f\n", n.Pair, n.Price, n.Quantity, n.VWAP)
		}

	}
}
