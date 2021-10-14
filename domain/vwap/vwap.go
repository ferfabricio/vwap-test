package vwap

import "errors"

type CalculationUnit struct {
	TotalPrice    float32
	TotalQuantity float32
	Result        float32
}

type Vwap struct {
	Pairs map[string]CalculationUnit
}

func (v Vwap) AddPair(key string) {
	v.Pairs[key] = CalculationUnit{
		TotalPrice:    0.0,
		TotalQuantity: 0.0,
		Result:        0.0,
	}
}

func (v Vwap) AddTrade(key string, price float32, quantity float32) error {
	p, ok := v.Pairs[key]
	if !ok {
		return errors.New("initilize the pair before add trades")
	}
	p.TotalPrice += price
	p.TotalQuantity += quantity
	v.Pairs[key] = p
	return nil
}
