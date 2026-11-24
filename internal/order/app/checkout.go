package app

import (
	"errors"
	"time"

	invdomain "example.com/shop/internal/inventory/domain"
	invinfra "example.com/shop/internal/inventory/infra"
	"example.com/shop/internal/order/domain"
	"example.com/shop/internal/order/infra"
)

const maxCheckoutAttempts = 3

type CheckoutService struct {
	orders  *infra.OrderRepository
	coupons map[string]*domain.Coupon
	stocks  *invinfra.StockRepository
}

func NewCheckoutService(orders *infra.OrderRepository, coupons map[string]*domain.Coupon, stocks *invinfra.StockRepository) *CheckoutService {
	return &CheckoutService{orders: orders, coupons: coupons, stocks: stocks}
}

// Receipt は注文確定時に購入者へ通知する内容。
type Receipt struct {
	OrderID string
	Lines   []domain.Line // 割引後の単価
	Total   int64
	// Remaining は引き当て後の在庫数（SKUごと）。残りわずか表示に使う。
	Remaining map[string]int
}

// Checkout はクーポンを適用し、在庫を引き当てて注文を確定する。
// 同時更新でバージョンが衝突した場合は読み直して再試行する。
func (s *CheckoutService) Checkout(orderID, couponCode string, now time.Time) (*Receipt, error) {
	var lastErr error
	for i := 0; i < maxCheckoutAttempts; i++ {
		r, err := s.checkoutOnce(orderID, couponCode, now)
		if err == nil {
			return r, nil
		}
		if errors.Is(err, domain.ErrVersionConflict) {
			lastErr = err
			continue
		}
		return nil, err
	}
	return nil, lastErr
}

func (s *CheckoutService) checkoutOnce(orderID, couponCode string, now time.Time) (*Receipt, error) {
	o, err := s.orders.Find(orderID)
	if err != nil {
		return nil, err
	}

	if couponCode != "" {
		c, ok := s.coupons[couponCode]
		if !ok {
			return nil, errors.New("unknown coupon")
		}
		if err := o.ApplyCoupon(c, now); err != nil {
			return nil, err
		}
	}

	remaining, err := s.reserve(o)
	if err != nil {
		return nil, err
	}

	receiptLines := o.Lines()
	if c := o.Coupon(); c != nil {
		for i := range receiptLines {
			receiptLines[i].UnitPrice -= c.DiscountFor(receiptLines[i].UnitPrice)
		}
	}

	if err := o.Place(); err != nil {
		return nil, err
	}
	if err := s.orders.Save(o); err != nil {
		return nil, err
	}

	return &Receipt{
		OrderID:   o.ID,
		Lines:     receiptLines,
		Total:     o.Total(),
		Remaining: remaining,
	}, nil
}

// reserve は注文の全明細の在庫を引き当て、引き当て後の在庫数を返す。
// 1つでも足りなければ何も引き当てない。
func (s *CheckoutService) reserve(o *domain.Order) (map[string]int, error) {
	s.stocks.Lock()
	defer s.stocks.Unlock()

	stocks := make([]*invdomain.Stock, 0, len(o.Lines()))
	for _, l := range o.Lines() {
		st, err := s.stocks.Find(l.SKU)
		if err != nil {
			return nil, err
		}
		if st.Available()+st.ReservedBy[o.ID] < l.Quantity {
			return nil, invdomain.ErrInsufficientStock
		}
		stocks = append(stocks, st)
	}

	remaining := make(map[string]int, len(stocks))
	for i, l := range o.Lines() {
		st := stocks[i]
		if err := st.Reserve(o.ID, l.Quantity); err != nil {
			return nil, err
		}
		s.stocks.Put(st)
		remaining[st.SKU] = st.Available()
	}
	return remaining, nil
}
