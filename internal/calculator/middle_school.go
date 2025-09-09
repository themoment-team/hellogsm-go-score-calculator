package calculator

import (
	"math/big"
	"hellogsm-go-score-calculator/internal/types"
)

func BuildCalcDtoWithFillEmpty(dto types.MiddleSchoolAchievementReqDto, graduationType types.GraduationType) types.MiddleSchoolAchievementCalcDto {
	if graduationType == types.GED {
		// GedAvgScore를 float64에서 big.Rat으로 변환
		var gedAvgScore *big.Rat
		if dto.GedAvgScore != nil {
			gedAvgScore = new(big.Rat).SetFloat64(*dto.GedAvgScore)
		}
		return types.MiddleSchoolAchievementCalcDto{
			GedAvgScore: gedAvgScore,
		}
	}
	
	// 졸업예정자 & 졸업자는 없는 성적을 복사하여 사용
	// 자유학년제(1학년)을 제외하고, 두개 이상의 빈학기가 없다는 가정
	tmpAchievement1_1 := dto.Achievement1_1
	tmpAchievement1_2 := dto.Achievement1_2
	tmpAchievement2_1 := dto.Achievement2_1
	tmpAchievement2_2 := dto.Achievement2_2
	tmpAchievement3_1 := dto.Achievement3_1
	tmpAchievement3_2 := dto.Achievement3_2

	// Spring 로직과 정확히 동일하게 수정 (else if 체인 제거)
	if graduationType == types.GRADUATE && tmpAchievement3_2 == nil {
		tmpAchievement3_2 = tmpAchievement3_1
	}
	if tmpAchievement3_1 == nil {
		tmpAchievement3_1 = tmpAchievement3_2
	}
	if tmpAchievement2_1 == nil {
		tmpAchievement2_1 = tmpAchievement2_2
	}
	if tmpAchievement2_2 == nil {
		tmpAchievement2_2 = tmpAchievement2_1
	}
	if graduationType == types.CANDIDATE && tmpAchievement1_2 == nil {
		if tmpAchievement1_1 == nil {
			tmpAchievement1_2 = tmpAchievement2_2
		} else {
			tmpAchievement1_2 = tmpAchievement1_1
		}
	}

	return types.MiddleSchoolAchievementCalcDto{
		Achievement1_2:          tmpAchievement1_2,
		Achievement2_1:          tmpAchievement2_1,
		Achievement2_2:          tmpAchievement2_2,
		Achievement3_1:          tmpAchievement3_1,
		Achievement3_2:          tmpAchievement3_2,
		ArtsPhysicalAchievement: dto.ArtsPhysicalAchievement,
		AbsentDays:              dto.AbsentDays,
		AttendanceDays:          dto.AttendanceDays,
		VolunteerTime:           dto.VolunteerTime,
		LiberalSystem:           dto.LiberalSystem,
		FreeSemester:            dto.FreeSemester,
		GedAvgScore:             nil, // 비 GED 경우는 사용하지 않음
	}
}

func CalcGeneralSubjectsSemesterScore(dto types.MiddleSchoolAchievementCalcDto, graduationType types.GraduationType) types.GeneralSubjectsSemesterScoreCalcDto {
	switch graduationType {
	case types.CANDIDATE:
		return types.GeneralSubjectsSemesterScoreCalcDto{
			Score1_2: CalcGeneralSubjectsScore(dto.Achievement1_2, big.NewRat(18, 1)),
			Score2_1: CalcGeneralSubjectsScore(dto.Achievement2_1, big.NewRat(45, 1)),
			Score2_2: CalcGeneralSubjectsScore(dto.Achievement2_2, big.NewRat(45, 1)),
			Score3_1: CalcGeneralSubjectsScore(dto.Achievement3_1, big.NewRat(72, 1)),
			Score3_2: big.NewRat(0, 1),
		}
	case types.GRADUATE:
		return types.GeneralSubjectsSemesterScoreCalcDto{
			Score1_2: big.NewRat(0, 1),
			Score2_1: CalcGeneralSubjectsScore(dto.Achievement2_1, big.NewRat(36, 1)),
			Score2_2: CalcGeneralSubjectsScore(dto.Achievement2_2, big.NewRat(36, 1)),
			Score3_1: CalcGeneralSubjectsScore(dto.Achievement3_1, big.NewRat(54, 1)),
			Score3_2: CalcGeneralSubjectsScore(dto.Achievement3_2, big.NewRat(54, 1)),
		}
	default:
		return types.GeneralSubjectsSemesterScoreCalcDto{}
	}
}

func CalcGeneralSubjectsTotalScore(generalSubjectsSemesterScore types.GeneralSubjectsSemesterScoreCalcDto) *big.Rat {
	// Stream.of()와 동일한 방식으로 처리
	scores := []*big.Rat{
		generalSubjectsSemesterScore.Score1_2,
		generalSubjectsSemesterScore.Score2_1,
		generalSubjectsSemesterScore.Score2_2,
		generalSubjectsSemesterScore.Score3_1,
		generalSubjectsSemesterScore.Score3_2,
	}

	total := big.NewRat(0, 1)
	for _, score := range scores {
		if score != nil {
			total.Add(total, score)
		}
	}

	return RoundToThreeDecimals(total)
}

func CalcGeneralSubjectsScore(achievements []int, maxPoint *big.Rat) *big.Rat {
	// 해당 학기의 등급별 점수 배열이 비어있거나 해당 학기의 배점이 없다면 0.000을 반환
	if achievements == nil || len(achievements) == 0 || maxPoint.Cmp(big.NewRat(0, 1)) == 0 {
		return RoundToThreeDecimals(big.NewRat(0, 1))
	}

	// 해당 학기에 수강하지 않은 과목이 있다면 제거한 리스트를 반환 (점수가 0인 원소 제거)
	var noZeroAchievements []int
	totalSum := 0
	for _, achievement := range achievements {
		totalSum += achievement
		if achievement != 0 {
			noZeroAchievements = append(noZeroAchievements, achievement)
		}
	}
	
	// 위에서 구한 리스트가 비어있다면 0.000을 반환
	if len(noZeroAchievements) == 0 {
		return RoundToThreeDecimals(big.NewRat(0, 1))
	}

	// 1. 점수로 환산된 등급을 모두 더한다.
	totalSumRat := big.NewRat(int64(totalSum), 1)
	
	// 2. 더한값 / (과목 수 * 5) (소수점 6째자리에서 반올림)
	divisor := big.NewRat(int64(len(noZeroAchievements)*5), 1)
	divideResult := new(big.Rat).Quo(totalSumRat, divisor)
	
	// 3. 각 학기당 배점 * 나눈값 (소수점 4째자리에서 반올림)
	result := new(big.Rat).Mul(divideResult, maxPoint)

	return RoundToThreeDecimals(result)
}

func CalcArtSportsScore(achievements []int) *big.Rat {
	if achievements == nil || len(achievements) == 0 {
		return RoundToThreeDecimals(big.NewRat(0, 1))
	}

	// 1. 각 등급별 갯수에 등급별 배점을 곱한 값을 더한다.
	totalScores := 0
	// 2. 각 등급별 갯수를 모두 더해 성취 수를 구한다.
	achievementCount := 0

	for _, achievement := range achievements {
		totalScores += achievement
		if achievement != 0 {
			achievementCount++
		}
	}

	// 과목 수가 0일시 0점 반환
	if achievementCount == 0 {
		return RoundToThreeDecimals(big.NewRat(0, 1))
	}

	// 3. 각 등급별 갯수를 더한 값(성취 수)에 5를 곱해 총점을 구한다.
	maxScore := 5 * achievementCount

	averageOfAchievementScale := big.NewRat(int64(totalScores), int64(maxScore))
	result := new(big.Rat).Mul(big.NewRat(60, 1), averageOfAchievementScale)

	return RoundToThreeDecimals(result)
}

func CalcAttendanceScore(absentDays, attendanceDays []int) *big.Rat {
	// 결석 횟수를 더함
	totalAbsentDays := 0
	for _, day := range absentDays {
		totalAbsentDays += day
	}

	// 결석 횟수가 10회 이상 0점을 반환
	if totalAbsentDays >= 10 {
		return big.NewRat(0, 1)
	}

	// 1. 모든 지각, 조퇴, 결과 횟수를 더함
	totalAttendanceDays := 0
	for _, day := range attendanceDays {
		totalAttendanceDays += day
	}

	// 2. 지각, 조퇴, 결과 횟수는 3개당 결석 1회
	absentResult := totalAttendanceDays / 3
	
	// 3. 총점(30점) - (3 * 총 결석 횟수)
	totalAbsent := totalAbsentDays + absentResult
	result := 30 - (3 * totalAbsent)

	// 총 점수가 0점 이하라면 0점을 반환
	if result <= 0 {
		return big.NewRat(0, 1)
	}

	return big.NewRat(int64(result), 1)
}

func CalcVolunteerScore(volunteerHours []int) *big.Rat {
	total := big.NewRat(0, 1)

	for _, hour := range volunteerHours {
		var score int64
		// 연간 7시간 이상
		if hour >= 7 {
			score = 10
		// 연간 6시간 이상
		} else if hour >= 6 {
			score = 8
		// 연간 5시간 이상
		} else if hour >= 5 {
			score = 6
		// 연간 4시간 이상
		} else if hour >= 4 {
			score = 4
		// 연간 3시간 이하
		} else {
			score = 2
		}
		total.Add(total, big.NewRat(score, 1))
	}

	return total
}
