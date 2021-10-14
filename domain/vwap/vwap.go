package vwap

import (
	"container/list"
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
	DataPoints             list.List
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
		DataPoints:             *list.New(),
	}
}

func calculateVwapResult(tp float32, tq float32) float32 {
	return tp / tq
}

func calculateTotalsInPair(p CalculationUnit, dp DataPoint, l int) CalculationUnit {
	if p.DataPoints.Len() >= l {
		f := p.DataPoints.Front()
		fv := f.Value.(DataPoint)
		p.TotalPricePlusQuantity -= (fv.Price * fv.Quantity)
		p.TotalQuantity -= fv.Quantity
		p.DataPoints.Remove(f)
	}
	p.DataPoints.PushBack(dp)
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

	v.Pairs[key] = calculateTotalsInPair(p, dp, v.Length)
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
