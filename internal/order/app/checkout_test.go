package app

import (
	"errors"
	"testing"

	invapp "example.com/shop/internal/inventory/app"
	invdomain "example.com/shop/internal/inventory/domain"
	invinfra "example.com/shop/internal/inventory/infra"
	"example.com/shop/internal/order/domain"
	"example.com/shop/internal/order/infra"
)

func setup(t *testing.T) (*CheckoutService, *infra.OrderRepository, *invinfra.StockRepository) {
	t.Helper()
	orders := infra.NewOrderRepository()
	stocks := invinfra.NewStockRepository()
	stocks.Lock()
	stocks.Put(invdomain.NewStock("A", 10))
	stocks.Put(invdomain.NewStock("B", 1))
	stocks.Unlock()
	return NewCheckoutService(orders, invapp.NewService(stocks)), orders, stocks
}

func TestCheckout(t *testing.T) {
	svc, orders, stocks := setup(t)
	o := domain.NewOrder("o1")
	o.AddLine("A", 2, 1000)
	if err := orders.Save(o); err != nil {
		t.Fatal(err)
	}

	if err := svc.Checkout("o1"); err != nil {
		t.Fatal(err)
	}

	got, _ := orders.Find("o1")
	if got.Status != domain.StatusPlaced {
		t.Errorf("status = %s", got.Status)
	}
	stocks.Lock()
	defer stocks.Unlock()
	st, _ := stocks.Find("A")
	if st.Available() != 8 {
		t.Errorf("available = %d", st.Available())
	}
}

func TestCheckoutInsufficientStock(t *testing.T) {
	svc, orders, _ := setup(t)
	o := domain.NewOrder("o1")
	o.AddLine("A", 1, 1000)
	o.AddLine("B", 2, 500)
	_ = orders.Save(o)

	err := svc.Checkout("o1")
	if !errors.Is(err, invdomain.ErrInsufficientStock) {
		t.Fatalf("err = %v", err)
	}
}
