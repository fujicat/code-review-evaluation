package domain

import "errors"

var (
	ErrNotFound        = errors.New("order not found")
	ErrVersionConflict = errors.New("order version conflict")
	ErrAlreadyPlaced   = errors.New("order already placed")
	ErrEmptyOrder      = errors.New("order has no lines")
	ErrCouponExpired   = errors.New("coupon expired")
)
