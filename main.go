package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	
	"hellogsm-go-score-calculator/internal/types"
	"hellogsm-go-score-calculator/internal/calculator"
	"hellogsm-go-score-calculator/internal/validator"
)

func createErrorResponse(statusCode int, errorCode, message string) events.APIGatewayProxyResponse {
	errorResp := types.ErrorResponse{
		Error:   "Validation Error",
		Message: message,
		Code:    errorCode,
	}
	
	body, _ := json.Marshal(errorResp)
	
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 디버깅을 위한 로그 추가
	log.Printf("Received request body: %s", request.Body)
	
	// 빈 body 체크
	if request.Body == "" {
		return createErrorResponse(400, "EMPTY_BODY", "Request body is empty"), nil
	}
	
	var reqDto types.MiddleSchoolAchievementReqDto
	if err := json.Unmarshal([]byte(request.Body), &reqDto); err != nil {
		log.Printf("JSON unmarshal error: %v", err)
		return createErrorResponse(400, "INVALID_JSON", fmt.Sprintf("Invalid JSON format: %s", err.Error())), nil
	}

	// 요청 데이터 검증
	if err := validator.ValidateRequest(reqDto); err != nil {
		log.Printf("Validation error: %v", err)
		return createErrorResponse(400, "VALIDATION_ERROR", err.Error()), nil
	}

	graduationType := types.GraduationType(reqDto.GraduationType)
	log.Printf("Processing with graduationType: %s", graduationType)
	
	calcDto := calculator.BuildCalcDtoWithFillEmpty(reqDto, graduationType)
	result := calculator.Execute(calcDto, graduationType)

	responseBody, err := json.Marshal(result)
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return createErrorResponse(500, "MARSHAL_ERROR", "Failed to create response"), nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBody),
	}, nil
}

func main() {
	lambda.Start(handler)
}
