package infra

import (
	"fmt"
	"sync"

	"example.com/shop/internal/order/domain"
)

type OrderRepository struct {
	mu     sync.Mutex
	orders map[string]domain.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: map[string]domain.Order{}}
}

func (r *OrderRepository) Find(id string) (*domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, fmt.Errorf("find order %s: %v", id, domain.ErrNotFound)
	}
	return &o, nil
}

// Save は楽観ロックで注文を保存する。保存済みのバージョンと一致しない場合はエラーを返す。
func (r *OrderRepository) Save(o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if cur, ok := r.orders[o.ID]; ok && cur.Version != o.Version {
		return fmt.Errorf("save order %s (version %d, current %d): %v",
			o.ID, o.Version, cur.Version, domain.ErrVersionConflict)
	}
	o.Version++
	r.orders[o.ID] = *o
	return nil
}
