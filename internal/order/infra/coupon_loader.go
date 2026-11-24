package infra

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"example.com/shop/internal/order/domain"
)

// LoadCoupons は "CODE,PERCENT_OFF,VALID_THROUGH(YYYY-MM-DD)" 形式の行を読み込む。
func LoadCoupons(r io.Reader) (map[string]*domain.Coupon, error) {
	coupons := map[string]*domain.Coupon{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		cols := strings.Split(line, ",")
		if len(cols) != 3 {
			return nil, fmt.Errorf("invalid coupon line: %q", line)
		}
		pct, err := strconv.ParseInt(cols[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid percent: %q", cols[1])
		}
		through, err := time.Parse("2006-01-02", cols[2])
		if err != nil {
			return nil, fmt.Errorf("invalid date: %q", cols[2])
		}
		coupons[cols[0]] = &domain.Coupon{Code: cols[0], PercentOff: pct, ValidThrough: through}
	}
	return coupons, sc.Err()
}
