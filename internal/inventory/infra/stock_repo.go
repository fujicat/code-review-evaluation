package infra

import (
	"fmt"
	"sync"

	"example.com/shop/internal/inventory/domain"
)

type StockRepository struct {
	mu     sync.Mutex
	stocks map[string]*domain.Stock
}

func NewStockRepository() *StockRepository {
	return &StockRepository{stocks: map[string]*domain.Stock{}}
}

func (r *StockRepository) Lock()   { r.mu.Lock() }
func (r *StockRepository) Unlock() { r.mu.Unlock() }

// Find は呼び出し側で Lock 済みであることを前提とする。
func (r *StockRepository) Find(sku string) (*domain.Stock, error) {
	s, ok := r.stocks[sku]
	if !ok {
		return nil, fmt.Errorf("stock %s not found", sku)
	}
	return s, nil
}

// Put は呼び出し側で Lock 済みであることを前提とする。
func (r *StockRepository) Put(s *domain.Stock) {
	r.stocks[s.SKU] = s
}
