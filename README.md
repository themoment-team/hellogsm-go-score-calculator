# hellogsm-go-score-calculator

## 개요
중학교 성적 계산 로직을 담당하는 Go 기반 Lambda 서비스입니다.
GSM(광주소프트웨어마이스터고) 입학 전형의 내신 성적 계산을 위한 핵심 엔진입니다.

## 주요 기능
- **일반고 졸업예정자/졸업자**: 중학교 학기별 성적 기반 계산
- **검정고시 출신**: 검정고시 평균 점수 기반 계산
- **성적 검증**: 입력 데이터 유효성 검사
- **고정밀도 계산**: big.Rat을 사용한 정확한 소수점 계산

## 계산 항목
### 교과 성적 (240점)
- **일반 교과**: 학기별 성적 환산 (최대 180점)
- **예체능 교과**: 예술/체육 성적 환산 (60점)

### 비교과 성적 (60점)  
- **출결 점수**: 결석/지각/조퇴 기반 계산 (30점)
- **봉사 점수**: 연간 봉사시간 기반 계산 (30점)

### 검정고시 특별 계산
- 평균 점수 기반 교과/비교과 점수 환산
- 출결 점수 만점(30점) 자동 부여

## 문제 해결
### 원본 문제
Lambda 테스트에서 다음과 같은 에러가 발생했습니다:
```json
{
  "statusCode": 400,
  "body": "{\"error\":\"Validation Error\",\"message\":\"Invalid JSON format: json: cannot unmarshal number into Go struct field MiddleSchoolAchievementReqDto.gedAvgScore of type *big.Rat\",\"code\":\"INVALID_JSON\"}"
}
```

테스트 입력:
```json
{
  "body": "{\"gedAvgScore\":92,\"graduationType\":\"GED\"}"
}
```

### 해결 방법
1. **MiddleSchoolAchievementReqDto**: `GedAvgScore *float64` (JSON 언마샬링용)
2. **MiddleSchoolAchievementCalcDto**: `GedAvgScore *big.Rat` (내부 계산용)
3. **BuildCalcDtoWithFillEmpty**: float64를 big.Rat으로 변환하는 로직 추가
4. **Validator**: float64 포인터를 검증하도록 수정

### 수정된 파일들
- `internal/types/types.go`: 구조체 필드 타입 변경
- `internal/calculator/middle_school.go`: BuildCalcDtoWithFillEmpty 함수에서 float64→big.Rat 변환
- `internal/validator/validator.go`: GedAvgScore 필드 검증 로직 수정
- `internal/calculator/ged.go`: 기존 big.Rat 계산 로직 유지

## 빌드 방법
```bash
cd /Users/snowykte0426/Programming/hellogsm-go-score-calculator
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bootstrap main.go
zip function.zip bootstrap
```

## 테스트
이제 다음 JSON이 정상적으로 처리됩니다:
```json
{
  "gedAvgScore": 92,
  "graduationType": "GED"
}
```

## 동작 원리
1. JSON에서 `gedAvgScore: 92`가 `*float64`로 언마샬링됨
2. `BuildCalcDtoWithFillEmpty`에서 `92.0`을 `big.Rat`으로 변환
3. 기존 GED 계산 로직이 그대로 동작하여 고정밀도 계산 수행
4. 최종 결과는 float64로 변환되어 응답

## 특징
- 기존 계산 로직은 그대로 유지 (big.Rat 사용)
- JSON 호환성 확보 (float64 사용)
- 타입 안전성 보장
- 고정밀도 계산 유지
