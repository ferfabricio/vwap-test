package vwap

import (
	"errors"
)

const defaultLength = 200

type VWAPConfig struct {
	Length int
}

type DataPoint struct {
	Price    float64
	Quantity float64
}
type CalculationUnit struct {
	TotalPricePlusQuantity float64
	TotalQuantity          float64
	Result                 float64
	DataPoints             []DataPoint
}

type Vwap struct {
	Pairs               map[string]CalculationUnit
	Length              int
	NotificationChannel chan CalculationEvent
}

type CalculationEvent struct {
	Pair     string
	Price    float64
	Quantity float64
	VWAP     float64
}

func (v Vwap) AddPair(key string) {
	v.Pairs[key] = CalculationUnit{
		TotalPricePlusQuantity: 0.0,
		TotalQuantity:          0.0,
		Result:                 0.0,
		DataPoints:             []DataPoint{},
	}
}

func calculateVwapResult(tp float64, tq float64) float64 {
	return tp / tq
}

func calculateTotalsInPair(p *CalculationUnit, dp DataPoint, l int) *CalculationUnit {
	if len(p.DataPoints) == l {
		// Remove first item from the array and the related value from the Calculation Unit
		f, pda := p.DataPoints[0], p.DataPoints[1:]
		p.TotalPricePlusQuantity -= (f.Price * f.Quantity)
		p.TotalQuantity -= f.Quantity
		p.DataPoints = pda
	}
	// Add the actual pair to the collection and update the calculation
	p.DataPoints = append(p.DataPoints, dp)
	p.TotalPricePlusQuantity += (dp.Price * dp.Quantity)
	p.TotalQuantity += dp.Quantity
	p.Result = calculateVwapResult(p.TotalPricePlusQuantity, p.TotalQuantity)

	return p
}

func (v Vwap) AddTrade(key string, price float64, quantity float64) error {
	p, ok := v.Pairs[key]
	if !ok {
		return errors.New("initilize the pair before add trades")
	}

	dp := DataPoint{
		Price:    price,
		Quantity: quantity,
	}

	// Update the pair calculation
	v.Pairs[key] = *calculateTotalsInPair(&p, dp, v.Length)

	// If the notification channel it is configured send a message
	if v.NotificationChannel != nil {
		v.NotificationChannel <- CalculationEvent{
			Pair:     key,
			Price:    price,
			Quantity: quantity,
			VWAP:     v.Pairs[key].Result,
		}
	}

	return nil
}

// Get VWAP calculation result
func (v Vwap) GetResult(key string) (float64, error) {
	p, ok := v.Pairs[key]
	if !ok {
		return 0, errors.New("pair not present")
	}

	return p.Result, nil
}

// Create new VWAP instance
func New() *Vwap {
	c := VWAPConfig{
		Length: defaultLength,
	}
	return &Vwap{
		Pairs:  map[string]CalculationUnit{},
		Length: c.Length,
	}
}

// Helper function used in tests
func NewWithNotification(en chan CalculationEvent) *Vwap {
	c := VWAPConfig{
		Length: defaultLength,
	}
	return &Vwap{
		Pairs:               map[string]CalculationUnit{},
		Length:              c.Length,
		NotificationChannel: en,
	}
}
