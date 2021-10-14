package vwap

import (
	"testing"

	"github.com/ferfabricio/vwap-calculation-engine/domain/currencies"
)

func TestAddPair(t *testing.T) {
	v := New()

	v.AddPair(currencies.BitcoinDollar)
	v.AddPair(currencies.EthereumBitcoin)
	v.AddPair(currencies.EthereumDollar)

	if len(v.Pairs) != 3 {
		t.Fail()
	}
}

func TestAddTrade(t *testing.T) {
	t.Run("Adding trade with success", func(t *testing.T) {
		v := New()
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
		v := New()
		err := v.AddTrade(currencies.BitcoinDollar, 10.1, 1.2)
		if err.Error() != "initilize the pair before add trades" {
			t.Error(err)
		}
	})
}

func TestVWAPCalcalationResult(t *testing.T) {
	sc := []struct {
		length    int
		dps       []DataPoint
		expResult float32
	}{
		{
			length: 2,
			dps: []DataPoint{
				{Price: 3438.23, Quantity: 0.05},
				{Price: 3438.16, Quantity: 0.06},
			},
			expResult: 3438.191818,
		},
		{
			length: 5,
			dps: []DataPoint{
				{Price: 3438.23, Quantity: 0.05},
				{Price: 3438.16, Quantity: 0.06},
				{Price: 3438.22, Quantity: 0.07},
				{Price: 3438.22, Quantity: 0.08},
				{Price: 3438.26, Quantity: 0.09},
			},
			expResult: 3438.221429,
		},
		{
			length: 2,
			dps: []DataPoint{
				{Price: 3438.23, Quantity: 0.05},
				{Price: 3438.16, Quantity: 0.06},
				{Price: 3438.22, Quantity: 0.07},
				{Price: 3438.22, Quantity: 0.08},
				{Price: 3438.26, Quantity: 0.09},
			},
			expResult: 3438.241699,
		},
	}
	for _, s := range sc {
		v := Vwap{
			Pairs:  map[string]CalculationUnit{},
			Length: s.length,
		}
		v.AddPair(currencies.BitcoinDollar)
		for _, p := range s.dps {
			err := v.AddTrade(currencies.BitcoinDollar, p.Price, p.Quantity)
			if err != nil {
				t.Fail()
			}
		}
		if vf, _ := v.GetResult(currencies.BitcoinDollar); vf != s.expResult {
			t.Errorf("expected result %f, is different from calculated %f", s.expResult, vf)
		}
	}
}
