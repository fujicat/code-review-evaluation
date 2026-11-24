package app

import (
	invapp "example.com/shop/internal/inventory/app"
	"example.com/shop/internal/order/infra"
)

type CheckoutService struct {
	orders    *infra.OrderRepository
	inventory *invapp.Service
}

func NewCheckoutService(orders *infra.OrderRepository, inventory *invapp.Service) *CheckoutService {
	return &CheckoutService{orders: orders, inventory: inventory}
}

// Checkout は在庫を引き当てて注文を確定する。
func (s *CheckoutService) Checkout(orderID string) error {
	o, err := s.orders.Find(orderID)
	if err != nil {
		return err
	}

	lines := make([]invapp.ReserveLine, 0, len(o.Lines()))
	for _, l := range o.Lines() {
		lines = append(lines, invapp.ReserveLine{SKU: l.SKU, Quantity: l.Quantity})
	}
	if err := s.inventory.ReserveForOrder(o.ID, lines); err != nil {
		return err
	}

	if err := o.Place(); err != nil {
		return err
	}
	return s.orders.Save(o)
}
