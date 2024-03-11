package helpers

import "math"

func ConvertAndRound(n float32, isStripe bool) int64 {
	if isStripe {
		return int64(math.Round(float64(n)*100) * 100)
	} else {
		return int64(math.Round(float64(n) * 100))
	}
}
