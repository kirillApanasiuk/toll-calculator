package main

import (
	"fmt"
	"github.com/kirillApanasiuk/toll-calculator/types"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(distance types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuID)
	if err != nil {
		return nil, fmt.Errorf("obu id (%d) not found", obuID)
	}
	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return inv, nil
}

func NewInvoiceAggregatore(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in the storage:", distance)
	return i.store.Insert(distance)
}
