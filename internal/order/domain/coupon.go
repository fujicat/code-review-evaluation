package domain

import "time"

type Coupon struct {
	Code       string
	PercentOff int64
	// ValidThrough はクーポンを利用できる最終日。日本時間でその日の終わり（23:59:59）まで有効。
	ValidThrough time.Time
}

// DiscountFor は金額 amount に対する割引額（円未満切り捨て）を返す。
func (c *Coupon) DiscountFor(amount int64) int64 {
	return amount * c.PercentOff / 100
}
