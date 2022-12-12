/*
 *  Copyright 2022 Nordeck IT + Consulting GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and  limitations
 *  under the License.
 *
 */

package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"feedback/internal/api"
	"feedback/internal/repository"
	gormjsonb "github.com/dariubs/gorm-jsonb"
	"github.com/golang-jwt/jwt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"
	"time"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Read(tokenValue string) (repository.Feedback, error) {
	feedback := repository.Feedback{
		BaseModel:     repository.BaseModel{uint(0), time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Rating:        3,
		RatingComment: "any_comment",
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	}
	args := m.Called(feedback)
	return repository.Feedback{}, args.Error(0)
}

func (m *RepositoryMock) Update(feedbackToUpdate repository.Feedback, tokenValue string) (repository.Feedback, error) {
	args := m.Called(feedbackToUpdate, tokenValue)
	return repository.Feedback{}, args.Error(0)
}

func (m *RepositoryMock) Store(value interface{}) error {
	// mock
	args := m.Called(value)
	return args.Error(0)

}

func Test_ValidTokenToJwt(t *testing.T) {
	repoMock := new(RepositoryMock)

	httpmock.Activate()
	response := httpmock.NewStringResponder(200, "{\n  \"results\": {\n    \"user\": true\n  },\n  \"user_id\": \"@user:domain.tld\"\n}")
	httpmock.RegisterResponder("POST", "https://some.url/verify/user", response)

	controller := New(repoMock, nil)
	request := httptest.NewRequest(http.MethodGet, "/token", nil)
	request.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJpbmNvbWluZyIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.0DoNIGqCNa1Tc41par4bzqnQWwlgsKCIP2mgUYEHemM")
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 200, responseWriter.Result().StatusCode)
	actual := responseWriter.Body.String()
	assert.True(t, strings.Contains(actual, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."), "not the same")
	repoMock.AssertExpectations(t)
}

func Test_ValidTokenToJwtWithOptions(t *testing.T) {
	repoMock := new(RepositoryMock)

	httpmock.Activate()
	response := httpmock.NewStringResponder(200, "{\n  \"results\": {\n    \"user\": true\n  },\n  \"user_id\": \"@user:domain.tld\"\n}")
	httpmock.RegisterResponder("POST", "https://some.url/verify/user", response)

	controller := New(repoMock, nil)
	request := httptest.NewRequest(http.MethodOptions, "/token", nil)
	request.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJpbmNvbWluZyIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.0DoNIGqCNa1Tc41par4bzqnQWwlgsKCIP2mgUYEHemM")
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 204, responseWriter.Result().StatusCode)
	actual := responseWriter.Body.String()
	assert.True(t, len(actual) == 0)
	repoMock.AssertExpectations(t)
}

func Test_InvalidResponse(t *testing.T) {
	repoMock := new(RepositoryMock)

	httpmock.Activate()
	response := httpmock.NewStringResponder(200, "{\n  \"results\": {\n    \"user\": false\n  },\n  \"user_id\": null\n}")
	httpmock.RegisterResponder("POST", "https://some.url/verify/user", response)

	controller := New(repoMock, nil)
	request := httptest.NewRequest(http.MethodGet, "/token", nil)
	request.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJpbmNvbWluZyIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.0DoNIGqCNa1Tc41par4bzqnQWwlgsKCIP2mgUYEHemM")
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 400, responseWriter.Result().StatusCode)
	actual := responseWriter.Body.String()
	assert.True(t, strings.Contains(actual, "user is not valid\n"), "not the same")
	repoMock.AssertExpectations(t)
}

func Test_InvalidToken(t *testing.T) {
	repoMock := new(RepositoryMock)

	httpmock.Activate()
	response := httpmock.NewStringResponder(200, "{\n  \"results\": {\n    \"user\": true\n  },\n  \"user_id\": \"@user:domain.tld\"\n}")
	httpmock.RegisterResponder("POST", "https://some.url/verify/user", response)

	controller := New(repoMock, nil)
	request := httptest.NewRequest(http.MethodGet, "/token", nil)
	request.Header.Set("authorization", "Bearer ")
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 400, responseWriter.Result().StatusCode)
	actual := responseWriter.Body.String()
	assert.True(t, strings.Contains(actual, "authentication header value has not matched / is not a bearer token\n"), "not the same")
	repoMock.AssertExpectations(t)
}

func TestController_CreateFeedback_Authorized(t *testing.T) {

	ratingComment := "any_comment"
	rating := 3

	repoMock := new(RepositoryMock)

	expected := &repository.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	}

	feedback := repository.Feedback{
		BaseModel:     repository.BaseModel{uint(0), time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Rating:        3,
		RatingComment: "any_comment",
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	}
	repoMock.On("Store", expected).Return(nil)
	repoMock.On("Read", feedback).Return(nil)
	controller := New(repoMock, nil)

	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	requestBody, _ := json.Marshal(&api.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      metadata,
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	})
	request := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(requestBody))
	mySigningKey := []byte("someArbitraryString")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{})
	signedTokenString, _ := token.SignedString(mySigningKey)
	request.Header.Set("authorization", "Bearer "+signedTokenString)
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 200, responseWriter.Result().StatusCode)
	repoMock.AssertExpectations(t)
	assert.Equal(t, "*", responseWriter.Result().Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "*", responseWriter.Result().Header.Get("Access-Control-Allow-Headers"))
}

func TestController_UpdateFeedback_Authorized(t *testing.T) {

	ratingComment := "any_comment"
	rating := 3

	repoMock := new(RepositoryMock)

	expected := &repository.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	}

	feedback := repository.Feedback{
		BaseModel:     repository.BaseModel{uint(0), time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Rating:        3,
		RatingComment: "any_comment",
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	}

	repoMock.On("Store", expected).Return(nil)
	repoMock.On("Read", feedback).Return(nil)
	controller := New(repoMock, nil)

	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	requestBody, _ := json.Marshal(&api.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      metadata,
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	})
	request := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(requestBody))
	mySigningKey := []byte("someArbitraryString")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{})
	signedTokenString, _ := token.SignedString(mySigningKey)
	request.Header.Set("authorization", "Bearer "+signedTokenString)
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 200, responseWriter.Result().StatusCode)
	repoMock.AssertExpectations(t)
	assert.Equal(t, "*", responseWriter.Result().Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "*", responseWriter.Result().Header.Get("Access-Control-Allow-Headers"))
}

func TestController_CreateFeedback_Unauthorized_WrongSigningKey(t *testing.T) {

	ratingComment := "any_comment"
	rating := 3

	repoMock := new(RepositoryMock)

	controller := New(repoMock, nil)

	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	requestBody, _ := json.Marshal(&api.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      metadata,
	})
	request := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(requestBody))
	mySigningKey := []byte("someInvalidSigningKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{})
	signedTokenString, _ := token.SignedString(mySigningKey)
	request.Header.Set("authorization", "Bearer "+signedTokenString)
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 401, responseWriter.Result().StatusCode)
	repoMock.AssertExpectations(t)
}

func TestController_CreateFeedback_MalformedToken(t *testing.T) {

	ratingComment := "any_comment"
	rating := 3

	repoMock := new(RepositoryMock)

	expected := &repository.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
	}
	repoMock.On("Store", expected).Return(nil)
	controller := New(repoMock, nil)

	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	requestBody, _ := json.Marshal(&api.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      metadata,
	})
	request := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(requestBody))
	mySigningKey := []byte("someArbitraryString")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{})
	signedTokenString, _ := token.SignedString(mySigningKey)
	request.Header.Set("authorization", "Bearer "+strings.Split(signedTokenString, ".")[0])
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	assert.Equal(t, 401, responseWriter.Result().StatusCode)
}

func TestController_CreateFeedbackWithOptions(t *testing.T) {
	ratingComment := "any_comment"
	rating := 3

	repoMock := new(RepositoryMock)
	controller := New(repoMock, nil)

	metadata := map[string]interface{}{
		"first_key":  "first_value",
		"second_key": "second_value",
	}

	requestBody, _ := json.Marshal(&api.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      metadata,
	})
	request := httptest.NewRequest(http.MethodOptions, "/feedback", bytes.NewReader(requestBody))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)
	assert.Equal(t, 204, responseWriter.Result().StatusCode)
	assert.Equal(t, "GET,HEAD,PUT,PATCH,POST,DELETE", responseWriter.Result().Header.Get("Access-Control-Allow-Methods"))
}

func TestController_CreateFeedback_emptyBody(t *testing.T) {
	repoMock := new(RepositoryMock)
	controller := New(repoMock, nil)

	request := httptest.NewRequest(http.MethodPost, "/feedback", iotest.ErrReader(errors.New("error")))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 401, status)
	repoMock.AssertNotCalled(t, "Store", mock.Anything)
}

func TestController_CreateFeedback_invalidJson(t *testing.T) {
	repoMock := new(RepositoryMock)
	controller := New(repoMock, nil)

	request := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader("broken"))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 401, status)
	repoMock.AssertNotCalled(t, "Store", mock.Anything)
}

func TestController_CreateFeedback_databaseError(t *testing.T) {
	repoMock := new(RepositoryMock)

	feedback := repository.Feedback{
		BaseModel:     repository.BaseModel{uint(0), time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)},
		Rating:        3,
		RatingComment: "any_comment",
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
		Jwt:           "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.Z-0V0WjFAQpqLLynDdrYLZIDxzPs-nCVHNxFutGeZIs",
	}
	repoMock.On("Read", feedback).Return(nil)
	repoMock.On("Store", mock.Anything).Return(errors.New("error"))

	controller := New(repoMock, nil)
	requestBody, _ := json.Marshal(&repository.Feedback{
		Rating:        1,
		RatingComment: "any",
	})
	request := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(requestBody))
	mySigningKey := []byte("someArbitraryString")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{})
	signedTokenString, _ := token.SignedString(mySigningKey)
	request.Header.Set("authorization", "Bearer "+signedTokenString)
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 500, status)
}
