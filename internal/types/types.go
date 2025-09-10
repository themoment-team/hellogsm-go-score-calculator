package types

import (
	"fmt"
	"math"
	"math/big"
)

type GraduationType string

const (
	CANDIDATE GraduationType = "CANDIDATE"
	GRADUATE  GraduationType = "GRADUATE"
	GED       GraduationType = "GED"
)

type MiddleSchoolAchievementReqDto struct {
	Achievement1_1          []int    `json:"achievement1_1"`
	Achievement1_2          []int    `json:"achievement1_2"`
	Achievement2_1          []int    `json:"achievement2_1"`
	Achievement2_2          []int    `json:"achievement2_2"`
	Achievement3_1          []int    `json:"achievement3_1"`
	Achievement3_2          []int    `json:"achievement3_2"`
	ArtsPhysicalAchievement []int    `json:"artsPhysicalAchievement"`
	AbsentDays              []int    `json:"absentDays"`
	AttendanceDays          []int    `json:"attendanceDays"`
	VolunteerTime           []int    `json:"volunteerTime"`
	LiberalSystem           string   `json:"liberalSystem"`
	FreeSemester            string   `json:"freeSemester"`
	GedAvgScore             *float64 `json:"gedAvgScore"`
	GraduationType          string   `json:"graduationType"`
	// 추가 필드들은 무시
	GeneralSubjects      []string `json:"generalSubjects"`
	ArtsPhysicalSubjects []string `json:"artsPhysicalSubjects"`
	NewSubjects          []string `json:"newSubjects"`
}

type MiddleSchoolAchievementCalcDto struct {
	Achievement1_2          []int    `json:"achievement1_2"`
	Achievement2_1          []int    `json:"achievement2_1"`
	Achievement2_2          []int    `json:"achievement2_2"`
	Achievement3_1          []int    `json:"achievement3_1"`
	Achievement3_2          []int    `json:"achievement3_2"`
	ArtsPhysicalAchievement []int    `json:"artsPhysicalAchievement"`
	AbsentDays              []int    `json:"absentDays"`
	AttendanceDays          []int    `json:"attendanceDays"`
	VolunteerTime           []int    `json:"volunteerTime"`
	LiberalSystem           string   `json:"liberalSystem"`
	FreeSemester            string   `json:"freeSemester"`
	GedAvgScore             *big.Rat `json:"-"` // JSON에서 제외, 내부 계산용
}

type GeneralSubjectsSemesterScoreCalcDto struct {
	Score1_2 *big.Rat `json:"score1_2"`
	Score2_1 *big.Rat `json:"score2_1"`
	Score2_2 *big.Rat `json:"score2_2"`
	Score3_1 *big.Rat `json:"score3_1"`
	Score3_2 *big.Rat `json:"score3_2"`
}

type GeneralSubjectsScoreDetailResDto struct {
	Score1_2 *ScoreValue `json:"score1_2,omitempty"`
	Score2_1 *ScoreValue `json:"score2_1,omitempty"`
	Score2_2 *ScoreValue `json:"score2_2,omitempty"`
	Score3_1 *ScoreValue `json:"score3_1,omitempty"`
	Score3_2 *ScoreValue `json:"score3_2,omitempty"`
}

type CalculatedScoreResDto struct {
	GeneralSubjectsScore       *ScoreValue                       `json:"generalSubjectsScore,omitempty"`
	GeneralSubjectsScoreDetail *GeneralSubjectsScoreDetailResDto `json:"generalSubjectsScoreDetail,omitempty"`
	ArtsPhysicalSubjectsScore  *ScoreValue                       `json:"artsPhysicalSubjectsScore,omitempty"`
	TotalSubjectsScore         *ScoreValue                       `json:"totalSubjectsScore,omitempty"`
	AttendanceScore            *ScoreValue                       `json:"attendanceScore"`
	VolunteerScore             *ScoreValue                       `json:"volunteerScore"`
	TotalScore                 *ScoreValue                       `json:"totalScore"`
}

type ScoreValue struct {
	Value float64
}

func (s ScoreValue) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%.3f", s.Value)), nil
}

func NewScoreValue(rat *big.Rat) *ScoreValue {
	if rat == nil {
		return &ScoreValue{Value: 0.0}
	}

	multiplier := big.NewRat(1000, 1)
	temp := new(big.Rat).Mul(rat, multiplier)

	floatVal, _ := temp.Float64()
	rounded := math.Round(floatVal) // 4자리에서 반올림

	return &ScoreValue{Value: rounded / 1000.0}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
