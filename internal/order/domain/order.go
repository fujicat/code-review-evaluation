package domain

import "time"

type Status string

const (
	StatusDraft  Status = "draft"
	StatusPlaced Status = "placed"
)

type Line struct {
	SKU       string
	Quantity  int
	UnitPrice int64
}

func (l Line) Subtotal() int64 {
	return l.UnitPrice * int64(l.Quantity)
}

type Order struct {
	ID      string
	Version int
	Status  Status
	lines   []Line
	coupon  *Coupon
}

func NewOrder(id string) *Order {
	return &Order{ID: id, Status: StatusDraft}
}

func (o *Order) AddLine(sku string, qty int, unitPrice int64) {
	o.lines = append(o.lines, Line{SKU: sku, Quantity: qty, UnitPrice: unitPrice})
}

// Lines は注文明細を返す。
func (o *Order) Lines() []Line {
	return o.lines
}

func (o *Order) Coupon() *Coupon {
	return o.coupon
}

// Total は明細の合計からクーポン割引を差し引いた金額を返す。
func (o *Order) Total() int64 {
	var sum int64
	for _, l := range o.lines {
		sum += l.Subtotal()
	}
	if o.coupon != nil {
		sum -= o.coupon.DiscountFor(sum)
	}
	return sum
}

// ApplyCoupon は注文にクーポンを適用する。有効期限を過ぎたクーポンは適用できない。
func (o *Order) ApplyCoupon(c *Coupon, now time.Time) error {
	if o.Status == StatusPlaced {
		return ErrAlreadyPlaced
	}
	if now.After(c.ValidThrough) {
		return ErrCouponExpired
	}
	o.coupon = c
	return nil
}

func (o *Order) Place() error {
	if o.Status == StatusPlaced {
		return ErrAlreadyPlaced
	}
	if len(o.lines) == 0 {
		return ErrEmptyOrder
	}
	o.Status = StatusPlaced
	return nil
}
