package types

import "math/big"

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
	GedAvgScore             *big.Rat `json:"-"`  // JSON에서 제외, 내부 계산용
}

type GeneralSubjectsSemesterScoreCalcDto struct {
	Score1_2 *big.Rat `json:"score1_2"`
	Score2_1 *big.Rat `json:"score2_1"`
	Score2_2 *big.Rat `json:"score2_2"`
	Score3_1 *big.Rat `json:"score3_1"`
	Score3_2 *big.Rat `json:"score3_2"`
}

type GeneralSubjectsScoreDetailResDto struct {
	Score1_2 float64 `json:"score1_2,omitempty"`
	Score2_1 float64 `json:"score2_1,omitempty"`
	Score2_2 float64 `json:"score2_2,omitempty"`
	Score3_1 float64 `json:"score3_1,omitempty"`
	Score3_2 float64 `json:"score3_2,omitempty"`
}

type CalculatedScoreResDto struct {
	GeneralSubjectsScore       float64                           `json:"generalSubjectsScore,omitempty"`
	GeneralSubjectsScoreDetail *GeneralSubjectsScoreDetailResDto `json:"generalSubjectsScoreDetail,omitempty"`
	ArtsPhysicalSubjectsScore  float64                           `json:"artsPhysicalSubjectsScore,omitempty"`
	TotalSubjectsScore         float64                           `json:"totalSubjectsScore,omitempty"`
	AttendanceScore            float64                           `json:"attendanceScore"`
	VolunteerScore             float64                           `json:"volunteerScore"`
	TotalScore                 float64                           `json:"totalScore"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
