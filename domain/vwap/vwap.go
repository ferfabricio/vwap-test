package vwap

import (
	"container/list"
	"errors"
)

type DataPoint struct {
	Price    float32
	Quantity float32
}
type CalculationUnit struct {
	TotalPrice    float32
	TotalQuantity float32
	Result        float32
	DataPoints    list.List
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

func calculateVwapResult(tp float32, tq float32, v float32) float32 {
	return ((tp / tq) * v) / tq
}

func calculateTotalsInPair(p CalculationUnit, dp DataPoint) CalculationUnit {
	if p.DataPoints.Len() >= 1 {
		f := p.DataPoints.Front()
		fv := f.Value.(DataPoint)
		p.TotalPrice -= fv.Price
		p.TotalQuantity -= fv.Quantity
		p.DataPoints.Remove(f)
	}
	p.DataPoints.PushBack(dp)
	p.TotalPrice += dp.Price
	p.TotalQuantity += dp.Quantity
	p.Result = calculateVwapResult(p.TotalPrice, p.TotalQuantity, dp.Quantity)
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

	v.Pairs[key] = calculateTotalsInPair(p, dp)
	return nil
}
