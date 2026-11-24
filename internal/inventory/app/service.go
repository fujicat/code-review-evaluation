package app

import (
	"example.com/shop/internal/inventory/infra"
)

type ReserveLine struct {
	SKU      string
	Quantity int
}

type Service struct {
	repo *infra.StockRepository
}

func NewService(repo *infra.StockRepository) *Service {
	return &Service{repo: repo}
}

// ReserveForOrder は注文の全明細の在庫を引き当てる。1つでも足りなければ何も引き当てない。
func (s *Service) ReserveForOrder(orderID string, lines []ReserveLine) error {
	s.repo.Lock()
	defer s.repo.Unlock()

	for _, l := range lines {
		st, err := s.repo.Find(l.SKU)
		if err != nil {
			return err
		}
		if err := checkAvailable(st.Available()+st.ReservedBy[orderID], l); err != nil {
			return err
		}
	}
	for _, l := range lines {
		st, _ := s.repo.Find(l.SKU)
		if err := st.Reserve(orderID, l.Quantity); err != nil {
			return err
		}
		s.repo.Put(st)
	}
	return nil
}
