package vwap

import (
	"testing"

	"github.com/ferfabricio/vwap-calculation-engine/domain/currencies"
)

func TestAddPair(t *testing.T) {
	v := Vwap{
		Pairs: map[string]CalculationUnit{},
	}

	v.AddPair(currencies.BitcoinDollar)
	v.AddPair(currencies.EthereumBitcoin)
	v.AddPair(currencies.EthereumDollar)

	if len(v.Pairs) != 3 {
		t.Fail()
	}
}

func TestAddTrade(t *testing.T) {
	t.Run("Adding trade with success", func(t *testing.T) {
		v := Vwap{
			Pairs: map[string]CalculationUnit{},
		}
		v.AddPair(currencies.BitcoinDollar)
		err := v.AddTrade(currencies.BitcoinDollar, 3438.23, 0.05)
		if err != nil {
			t.Error(err)
		}
		p := v.Pairs[currencies.BitcoinDollar]
		if p.TotalPricePlusQuantity != 171.9115 || p.TotalQuantity != 0.05 {
			t.Fail()
		}
	})

	t.Run("Adding trade not existent pair", func(t *testing.T) {
		v := Vwap{
			Pairs: map[string]CalculationUnit{},
		}
		err := v.AddTrade(currencies.BitcoinDollar, 10.1, 1.2)
		if err.Error() != "initilize the pair before add trades" {
			t.Error(err)
		}
	})
}
