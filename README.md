# hellogsm-go-score-calculator

## 개요
중학교 성적 계산 로직을 담당하는 Go 기반 서비스입니다.
GSM(광주소프트웨어마이스터고) 입학 전형의 내신 성적 계산을 위한 엔진입니다

## 사용법

### API 개요
이 서비스는 AWS Lambda 함수로 배포되어 REST API 형태로 동작합니다.

### 요청 형식
**HTTP Method**: POST  
**Content-Type**: application/json

### 공통 필드 설명

#### 전형 타입 (graduationType)
- `"CANDIDATE"`: 재학생 (1-3학년 성적 모두 필요)
- `"GRADUATE"`: 졸업생 (1-3학년 성적 모두 필요)  
- `"GED"`: 검정고시 (평균 점수만 필요)

#### 일반 교과 성적 필드
- `achievement1_1`: 1학년 1학기 성적 배열 (5단계: 1~5)
- `achievement1_2`: 1학년 2학기 성적 배열
- `achievement2_1`: 2학년 1학기 성적 배열
- `achievement2_2`: 2학년 2학기 성적 배열
- `achievement3_1`: 3학년 1학기 성적 배열
- `achievement3_2`: 3학년 2학기 성적 배열

#### 예체능 교과 성적 필드
- `artsPhysicalAchievement`: 예체능 성적 배열 (전 학기 통합)

#### 비교과 필드
- `absentDays`: 학기별 결석일수 배열
- `attendanceDays`: 학기별 출석일수 배열  
- `volunteerTime`: 연도별 봉사시간 배열

#### 시스템 필드
- `liberalSystem`: 자유학기제 적용 여부 ("FIRST" 또는 "SECOND")
- `freeSemester`: 자유학기 대상 학기 ("1-1", "1-2", "2-1", "2-2")

#### 검정고시 전용 필드
- `gedAvgScore`: 검정고시 평균 점수 (숫자)

### 사용 예시

#### 1. 재학생/졸업생 전형

```json
{
  "graduationType": "CANDIDATE",
  "achievement1_1": [3, 4, 3, 5, 4],
  "achievement1_2": [4, 4, 3, 5, 4],
  "achievement2_1": [4, 5, 4, 5, 5],
  "achievement2_2": [4, 5, 4, 5, 5],
  "achievement3_1": [5, 5, 4, 5, 5],
  "achievement3_2": [5, 5, 5, 5, 5],
  "artsPhysicalAchievement": [4, 4, 5, 3, 4, 4],
  "absentDays": [2, 1, 0, 3, 1, 2],
  "attendanceDays": [190, 195, 200, 185, 198, 192],
  "volunteerTime": [25, 30, 35],
  "liberalSystem": "FIRST",
  "freeSemester": "1-2"
}
```

#### 2. 검정고시 전형

```json
{
  "graduationType": "GED",
  "gedAvgScore": 92.5
}
```

### 응답 형식

#### 성공 응답 (HTTP 200)
```json
{
  "generalSubjectsScore": 165.75,
  "generalSubjectsScoreDetail": {
    "score1_2": 32.4,
    "score2_1": 36.0,
    "score2_2": 36.0,
    "score3_1": 33.75,
    "score3_2": 27.6
  },
  "artsPhysicalSubjectsScore": 48.0,
  "totalSubjectsScore": 213.75,
  "attendanceScore": 28.5,
  "volunteerScore": 30.0,
  "totalScore": 272.25
}
```

#### 오류 응답 (HTTP 400/500)
```json
{
  "error": "Validation Error",
  "message": "상세 오류 메시지",
  "code": "ERROR_CODE"
}
```

### 점수 체계

#### 총점: 300점
- **교과 성적**: 240점
  - 일반 교과: 180점 (학기별 최대 36점 × 5학기)
  - 예체능 교과: 60점
- **비교과 성적**: 60점
  - 출결 점수: 30점
  - 봉사 점수: 30점

#### 검정고시 특별 계산
- 교과 점수: `(평균점수 - 60) / 40 × 240`
- 비교과 점수: 출결 30점(만점) + 봉사 30점(만점)

### 오류 코드
- `EMPTY_BODY`: 요청 본문이 비어있음
- `INVALID_JSON`: JSON 형식 오류
- `VALIDATION_ERROR`: 데이터 검증 실패
- `MARSHAL_ERROR`: 응답 생성 실패

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

## 빌드

#### Mac/Linux 환경
```bash
chmod +x build.sh
./build.sh
```

#### Windows 환경
```cmd
#Git Bash
bash build.sh

#PowerShell
powershell -ExecutionPolicy Bypass -File build.ps1
```

#### Windows 환경
```cmd
# 의존성 정리
go mod tidy

# Lambda용 바이너리 빌드 (PowerShell)
$env:CGO_ENABLED=0; $env:GOOS="linux"; $env:GOARCH="amd64"; go build -o bootstrap main.go

# ZIP 패키지 생성 (PowerShell - Compress-Archive 사용)
Compress-Archive -Path bootstrap -DestinationPath function.zip -Force

# 또는 7-Zip 사용 (설치된 경우)
7z a function.zip bootstrap
```

### 빌드 결과
성공적으로 빌드되면 다음 파일이 생성됩니다:
- `bootstrap`: AWS Lambda 실행 파일
- `function.zip`: Lambda 배포용 ZIP 패키지
