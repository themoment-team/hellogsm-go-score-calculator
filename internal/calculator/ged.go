package calculator

import (
	"go-hellogsm-score-calculator/internal/types"
	"math/big"
)

// GED 검정고시 계산 함수들
func ExecuteGed(dto types.MiddleSchoolAchievementCalcDto) types.CalculatedScoreResDto {
	// 검정고시 평균 점수
	averageScore := dto.GedAvgScore
	if averageScore == nil {
		// 이 경우는 검증에서 걸러져야 함
		return types.CalculatedScoreResDto{}
	}

	// 검정고시 교과 성적 환산 점수 (총점: 240점)
	gedTotalSubjectsScore := CalcGedTotalSubjectsScore(averageScore)

	// 검정고시 봉사 성적 환산 점수 (총점: 30점)
	gedVolunteerScore := CalcGedVolunteerScore(averageScore)

	// 검정고시 출결 성적 점수 (총점: 30점)
	gedAttendanceScore := big.NewRat(30, 1)

	// 검정고시 비 교과 성적 환산 점수 (총점: 60점)
	gedTotalNonSubjectsScore := new(big.Rat).Add(gedVolunteerScore, gedAttendanceScore)
	gedTotalNonSubjectsScore = RoundToThreeDecimals(gedTotalNonSubjectsScore)

	// 검정고시 총 점수 (교과 성적 + 비교과 성적) (총점: 300점)
	totalScore := new(big.Rat).Add(gedTotalSubjectsScore, gedTotalNonSubjectsScore)
	totalScore = RoundToThreeDecimals(totalScore)

	return types.CalculatedScoreResDto{
		TotalSubjectsScore: RatToFloat64(gedTotalSubjectsScore),
		AttendanceScore:    RatToFloat64(gedAttendanceScore),
		VolunteerScore:     RatToFloat64(gedVolunteerScore),
		TotalScore:         RatToFloat64(totalScore),
	}
}

func CalcGedTotalSubjectsScore(averageScore *big.Rat) *big.Rat {
	// (평균점수 - 50) / 50 * 240
	temp := new(big.Rat).Sub(averageScore, big.NewRat(50, 1))
	temp.Quo(temp, big.NewRat(50, 1))
	result := new(big.Rat).Mul(temp, big.NewRat(240, 1))
	result = RoundToThreeDecimals(result)

	// 0보다 작으면 0을 반환
	if result.Cmp(big.NewRat(0, 1)) < 0 {
		return big.NewRat(0, 1)
	}

	return result
}

func CalcGedVolunteerScore(averageScore *big.Rat) *big.Rat {
	// (평균점수 - 40) / 60 * 30
	temp := new(big.Rat).Sub(averageScore, big.NewRat(40, 1))
	temp.Quo(temp, big.NewRat(60, 1))
	result := new(big.Rat).Mul(temp, big.NewRat(30, 1))
	result = RoundToThreeDecimals(result)

	// 0보다 작으면 0을 반환
	if result.Cmp(big.NewRat(0, 1)) < 0 {
		return big.NewRat(0, 1)
	}

	return result
}
