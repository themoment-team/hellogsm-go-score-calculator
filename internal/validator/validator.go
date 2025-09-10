package validator

import (
	"fmt"
	"hellogsm-go-score-calculator/internal/types"
)

func ValidateRequest(dto types.MiddleSchoolAchievementReqDto) error {
	// graduationType 필수 검증
	if dto.GraduationType == "" {
		return fmt.Errorf("graduationType is required. Must be one of: CANDIDATE, GRADUATE, GED")
	}

	graduationType := types.GraduationType(dto.GraduationType)
	if graduationType != types.CANDIDATE && graduationType != types.GRADUATE && graduationType != types.GED {
		return fmt.Errorf("invalid graduationType: %s. Must be one of: CANDIDATE, GRADUATE, GED", dto.GraduationType)
	}

	// GED인 경우 gedAvgScore 필수 검증
	if graduationType == types.GED {
		if dto.GedAvgScore == nil {
			return fmt.Errorf("gedAvgScore is required for GED graduationType")
		}

		gedScoreFloat := *dto.GedAvgScore
		if gedScoreFloat < 60 || gedScoreFloat > 100 {
			return fmt.Errorf("gedAvgScore must be between 60 and 100: %.3f", gedScoreFloat)
		}

		// GED는 다른 검증을 건너뛰기
		return nil
	}

	// liberalSystem 검증 (CANDIDATE, GRADUATE만)
	if dto.LiberalSystem == "" {
		return fmt.Errorf("liberalSystem is required")
	}
	if dto.LiberalSystem != "자유학기제" && dto.LiberalSystem != "자유학년제" {
		return fmt.Errorf("invalid liberalSystem: %s. Must be '자유학기제' or '자유학년제'", dto.LiberalSystem)
	}

	// 자유학기제인 경우 freeSemester 검증
	if dto.LiberalSystem == "자유학기제" && dto.FreeSemester == "" {
		return fmt.Errorf("freeSemester is required when liberalSystem is '자유학기제'")
	}

	validSemesters := []string{"", "1-1", "1-2", "2-1", "2-2", "3-1", "3-2"}
	if dto.LiberalSystem == "자유학기제" {
		found := false
		for _, semester := range validSemesters {
			if dto.FreeSemester == semester {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("invalid freeSemester: %s. Must be one of: 1-1, 1-2, 2-1, 2-2, 3-1, 3-2", dto.FreeSemester)
		}
	}

	if dto.ArtsPhysicalAchievement == nil || len(dto.ArtsPhysicalAchievement) == 0 {
		return fmt.Errorf("artsPhysicalAchievement is required and cannot be empty")
	}

	if dto.AbsentDays == nil || len(dto.AbsentDays) != 3 {
		return fmt.Errorf("absentDays must contain exactly 3 elements (for 3 years)")
	}

	if dto.AttendanceDays == nil || len(dto.AttendanceDays) != 9 {
		return fmt.Errorf("attendanceDays must contain exactly 9 elements")
	}

	if dto.VolunteerTime == nil || len(dto.VolunteerTime) != 3 {
		return fmt.Errorf("volunteerTime must contain exactly 3 elements (for 3 years)")
	}

	// 음수 값 검증
	for i, day := range dto.AbsentDays {
		if day < 0 {
			return fmt.Errorf("absentDays[%d] cannot be negative: %d", i, day)
		}
	}

	for i, day := range dto.AttendanceDays {
		if day < 0 {
			return fmt.Errorf("attendanceDays[%d] cannot be negative: %d", i, day)
		}
	}

	for i, hour := range dto.VolunteerTime {
		if hour < 0 {
			return fmt.Errorf("volunteerTime[%d] cannot be negative: %d", i, hour)
		}
	}

	// 일반교과 등급 검증
	validateGrades := func(grades []int, name string) error {
		if grades == nil {
			return nil // null은 허용
		}
		for i, grade := range grades {
			if grade != 0 && (grade < 1 || grade > 5) {
				return fmt.Errorf("%s[%d] must be between 1-5 or 0 (for no grade): %d", name, i, grade)
			}
		}
		return nil
	}

	if err := validateGrades(dto.Achievement1_1, "achievement1_1"); err != nil {
		return err
	}
	if err := validateGrades(dto.Achievement1_2, "achievement1_2"); err != nil {
		return err
	}
	if err := validateGrades(dto.Achievement2_1, "achievement2_1"); err != nil {
		return err
	}
	if err := validateGrades(dto.Achievement2_2, "achievement2_2"); err != nil {
		return err
	}
	if err := validateGrades(dto.Achievement3_1, "achievement3_1"); err != nil {
		return err
	}
	if err := validateGrades(dto.Achievement3_2, "achievement3_2"); err != nil {
		return err
	}

	// 예체능 등급 검증
	for i, grade := range dto.ArtsPhysicalAchievement {
		if grade != 0 && (grade < 3 || grade > 5) {
			return fmt.Errorf("artsPhysicalAchievement[%d] must be between 3-5 or 0 (for no grade): %d", i, grade)
		}
	}

	return nil
}
