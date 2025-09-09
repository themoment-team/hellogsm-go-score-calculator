package calculator

import (
	"math"
	"math/big"
)

// big.Rat을 소수점 3자리 float64로 변환
func RatToFloat64(rat *big.Rat) float64 {
	if rat == nil {
		return 0.0
	}

	multiplier := big.NewRat(1000, 1)
	temp := new(big.Rat).Mul(rat, multiplier)

	floatVal, _ := temp.Float64()
	rounded := math.Round(floatVal)

	return rounded / 1000.0
}

// 소수점 3자리에서 반올림하는 함수 (Java의 HALF_UP과 동일)
func RoundToThreeDecimals(value *big.Rat) *big.Rat {
	result := new(big.Rat)
	multiplier := big.NewRat(1000, 1)
	temp := new(big.Rat).Mul(value, multiplier)

	floatVal, _ := temp.Float64()
	rounded := math.Round(floatVal)

	result.SetInt64(int64(rounded))
	result.Quo(result, multiplier)
	return result
}
