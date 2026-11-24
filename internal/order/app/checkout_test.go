package app

import (
	"errors"
	"testing"
	"time"

	invdomain "example.com/shop/internal/inventory/domain"
	invinfra "example.com/shop/internal/inventory/infra"
	"example.com/shop/internal/order/domain"
	"example.com/shop/internal/order/infra"
)

var jst = time.FixedZone("JST", 9*60*60)

func setup(t *testing.T) (*CheckoutService, *infra.OrderRepository, *invinfra.StockRepository) {
	t.Helper()
	orders := infra.NewOrderRepository()
	stocks := invinfra.NewStockRepository()
	stocks.Lock()
	stocks.Put(invdomain.NewStock("A", 10))
	stocks.Put(invdomain.NewStock("B", 1))
	stocks.Unlock()
	coupons := map[string]*domain.Coupon{
		"WINTER10": {Code: "WINTER10", PercentOff: 10, ValidThrough: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)},
	}
	return NewCheckoutService(orders, coupons, stocks), orders, stocks
}

func TestCheckout(t *testing.T) {
	svc, orders, _ := setup(t)
	o := domain.NewOrder("o1")
	o.AddLine("A", 2, 1000)
	if err := orders.Save(o); err != nil {
		t.Fatal(err)
	}

	r, err := svc.Checkout("o1", "", time.Date(2026, 12, 1, 10, 0, 0, 0, jst))
	if err != nil {
		t.Fatal(err)
	}

	got, _ := orders.Find("o1")
	if got.Status != domain.StatusPlaced {
		t.Errorf("status = %s", got.Status)
	}
	if r.Remaining["A"] != 8 {
		t.Errorf("remaining = %d", r.Remaining["A"])
	}
}

func TestCheckoutWithCoupon(t *testing.T) {
	svc, orders, _ := setup(t)
	o := domain.NewOrder("o1")
	o.AddLine("A", 1, 1000)
	_ = orders.Save(o)

	r, err := svc.Checkout("o1", "WINTER10", time.Date(2026, 12, 1, 10, 0, 0, 0, jst))
	if err != nil {
		t.Fatal(err)
	}
	if r.Lines[0].UnitPrice != 900 {
		t.Errorf("unit price = %d", r.Lines[0].UnitPrice)
	}
}

func TestCheckoutExpiredCoupon(t *testing.T) {
	svc, orders, _ := setup(t)
	o := domain.NewOrder("o1")
	o.AddLine("A", 1, 1000)
	_ = orders.Save(o)

	_, err := svc.Checkout("o1", "WINTER10", time.Date(2027, 1, 1, 10, 0, 0, 0, jst))
	if !errors.Is(err, domain.ErrCouponExpired) {
		t.Fatalf("err = %v", err)
	}
}

func TestCheckoutInsufficientStock(t *testing.T) {
	svc, orders, _ := setup(t)
	o := domain.NewOrder("o1")
	o.AddLine("A", 1, 1000)
	o.AddLine("B", 2, 500)
	_ = orders.Save(o)

	_, err := svc.Checkout("o1", "", time.Date(2026, 12, 1, 10, 0, 0, 0, jst))
	if !errors.Is(err, invdomain.ErrInsufficientStock) {
		t.Fatalf("err = %v", err)
	}
}
