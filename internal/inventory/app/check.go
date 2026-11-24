package app

import (
	"fmt"

	"example.com/shop/internal/inventory/domain"
)

func checkAvailable(available int, l ReserveLine) error {
	if available < l.Quantity {
		return fmt.Errorf("sku %s: want %d, available %d: %w", l.SKU, l.Quantity, available, domain.ErrInsufficientStock)
	}
	return nil
}
