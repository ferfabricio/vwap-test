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
	err = c.Configure([]string{
		currencies.BitcoinDollar,
		currencies.EthereumBitcoin,
		currencies.EthereumDollar,
	})

	if err != nil {
		panic(err)
	}

	nf := make(chan vwap.CalculationEvent, 200)
	v := vwap.NewWithNotification(nf)
	v.AddPair(currencies.BitcoinDollar)
	v.AddPair(currencies.EthereumBitcoin)
	v.AddPair(currencies.EthereumDollar)

	ch := make(chan coinbase.MatchMessage, 200)
	go c.GetData(ch)
	defer close(ch)
	for {
		select {
		case r := <-ch:
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
			fmt.Printf("Pair: %s, Price: %f, Quantity: %f, VWAP: %f\n", n.Pair, n.Price, n.Quantity, n.VWAP)
		}

	}
}
