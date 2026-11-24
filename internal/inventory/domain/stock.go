package domain

import (
	"errors"
	"fmt"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type Stock struct {
	SKU        string
	OnHand     int
	Reserved   int
	ReservedBy map[string]int // orderID -> quantity
}

func NewStock(sku string, onHand int) *Stock {
	return &Stock{SKU: sku, OnHand: onHand, ReservedBy: map[string]int{}}
}

func (s *Stock) Available() int {
	return s.OnHand - s.Reserved
}

// Reserve は注文のために在庫を引き当てる。同じ注文での再引き当ては数量を上書きする。
func (s *Stock) Reserve(orderID string, qty int) error {
	prev := s.ReservedBy[orderID]
	if s.Available()+prev < qty {
		return fmt.Errorf("sku %s: want %d, available %d: %w", s.SKU, qty, s.Available()+prev, ErrInsufficientStock)
	}
	s.Reserved += qty - prev
	s.ReservedBy[orderID] = qty
	return nil
}
