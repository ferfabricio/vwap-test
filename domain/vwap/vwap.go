package vwap

import (
	"errors"
)

const defaultLength = 200

type VWAPConfig struct {
	Length int
}

type DataPoint struct {
	Price    float32
	Quantity float32
}
type CalculationUnit struct {
	TotalPricePlusQuantity float32
	TotalQuantity          float32
	Result                 float32
	DataPoints             []DataPoint
}

type Vwap struct {
	Pairs  map[string]CalculationUnit
	Length int
}

func (v Vwap) AddPair(key string) {
	v.Pairs[key] = CalculationUnit{
		TotalPricePlusQuantity: 0.0,
		TotalQuantity:          0.0,
		Result:                 0.0,
		DataPoints:             []DataPoint{},
	}
}

func calculateVwapResult(tp float32, tq float32) float32 {
	return tp / tq
}

func calculateTotalsInPair(p *CalculationUnit, dp DataPoint, l int) *CalculationUnit {
	if len(p.DataPoints) == l {
		f, pda := p.DataPoints[0], p.DataPoints[1:]
		p.TotalPricePlusQuantity -= (f.Price * f.Quantity)
		p.TotalQuantity -= f.Quantity
		p.DataPoints = pda
	}
	p.DataPoints = append(p.DataPoints, dp)
	p.TotalPricePlusQuantity += (dp.Price * dp.Quantity)
	p.TotalQuantity += dp.Quantity
	p.Result = calculateVwapResult(p.TotalPricePlusQuantity, p.TotalQuantity)
	return p
}

func (v Vwap) AddTrade(key string, price float32, quantity float32) error {
	p, ok := v.Pairs[key]
	if !ok {
		return errors.New("initilize the pair before add trades")
	}

	dp := DataPoint{
		Price:    price,
		Quantity: quantity,
	}

	v.Pairs[key] = *calculateTotalsInPair(&p, dp, v.Length)
	return nil
}

func (v Vwap) GetResult(key string) (float32, error) {
	p, ok := v.Pairs[key]
	if !ok {
		return 0, errors.New("pair not present")
	}

	return p.Result, nil
}

func New() *Vwap {
	c := VWAPConfig{
		Length: defaultLength,
	}
	return &Vwap{
		Pairs:  map[string]CalculationUnit{},
		Length: c.Length,
	}
}
