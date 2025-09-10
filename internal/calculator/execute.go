package calculator

import (
	"go-hellogsm-score-calculator/internal/types"
	"math/big"
)

func Execute(dto types.MiddleSchoolAchievementCalcDto, graduationType types.GraduationType) types.CalculatedScoreResDto {
	// GED 검정고시 처리
	if graduationType == types.GED {
		return ExecuteGed(dto)
	}

	generalSubjectsSemesterScore := CalcGeneralSubjectsSemesterScore(dto, graduationType)

	// 일반 교과 성적 환산값 (총점: 180점)
	generalSubjectsScore := CalcGeneralSubjectsTotalScore(generalSubjectsSemesterScore)

	// 예체능 성적 환산값 (총점: 60점)
	artsPhysicalSubjectsScore := CalcArtSportsScore(dto.ArtsPhysicalAchievement)

	// 교과 성적 환산값 (예체능 성적 + 일반 교과 성적, 총점: 240점)
	totalSubjectsScore := new(big.Rat).Add(artsPhysicalSubjectsScore, generalSubjectsScore)
	totalSubjectsScore = RoundToThreeDecimals(totalSubjectsScore)

	// 출결 성적 (총점: 30점)
	attendanceScore := CalcAttendanceScore(dto.AbsentDays, dto.AttendanceDays)
	attendanceScore = RoundToThreeDecimals(attendanceScore)

	// 봉사 성적 (총점: 30점)
	volunteerScore := CalcVolunteerScore(dto.VolunteerTime)
	volunteerScore = RoundToThreeDecimals(volunteerScore)

	// 비 교과 성적 환산값 (총점: 60점)
	totalNonSubjectsScore := new(big.Rat).Add(attendanceScore, volunteerScore)
	totalNonSubjectsScore = RoundToThreeDecimals(totalNonSubjectsScore)

	// 내신 성적 총 점수 (총점: 300점)
	totalScore := new(big.Rat).Add(totalSubjectsScore, totalNonSubjectsScore)
	totalScore = RoundToThreeDecimals(totalScore)

	// big.Rat을 ScoreValue로 변환
	generalSubjectsScoreDetail := &types.GeneralSubjectsScoreDetailResDto{
		Score1_2: types.NewScoreValue(generalSubjectsSemesterScore.Score1_2),
		Score2_1: types.NewScoreValue(generalSubjectsSemesterScore.Score2_1),
		Score2_2: types.NewScoreValue(generalSubjectsSemesterScore.Score2_2),
		Score3_1: types.NewScoreValue(generalSubjectsSemesterScore.Score3_1),
		Score3_2: types.NewScoreValue(generalSubjectsSemesterScore.Score3_2),
	}

	return types.CalculatedScoreResDto{
		GeneralSubjectsScore:       types.NewScoreValue(generalSubjectsScore),
		ArtsPhysicalSubjectsScore:  types.NewScoreValue(artsPhysicalSubjectsScore),
		AttendanceScore:            types.NewScoreValue(attendanceScore),
		VolunteerScore:             types.NewScoreValue(volunteerScore),
		TotalScore:                 types.NewScoreValue(totalScore),
		GeneralSubjectsScoreDetail: generalSubjectsScoreDetail,
	}
}
