package domain

import "math"

func PercentChange(oldPrice, newPrice float64) (float64, error) {
	if oldPrice <= 0 {
		return 0, ErrDivisionByZero
	}
	change := ((newPrice - oldPrice) / oldPrice) * 100
	return change, nil
}
func ReachedThreshold(changePercent, thresholdPercent float64) bool {
	if thresholdPercent <= 0 {
		return false
	}
	return math.Abs(changePercent) >= thresholdPercent
}
func ChangeDirection(changePercent float64) string {
	if changePercent > 0 {
		return "up"
	}
	if changePercent < 0 {
		return "down"
	}
	return "flat"
}
