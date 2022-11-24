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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) Store(value interface{}) error {
	args := m.Called(value)
	return args.Error(0)

}

func TestController_CreateToken(t *testing.T) {
	tokenValue := "{\n  \"acr\": \"1\",\n  \"aid\": \"default\",\n  \"amr\": [\n    \"pwd\"\n  ],\n  \"aud\": \"default-demo\",\n  \"auth_time\": 1631696786,\n  \"email\": \"\",\n  \"email_verified\": false,\n  \"exp\": 1631700395,\n  \"iat\": 1631696795,\n  \"idp\": \"default\",\n  \"iss\": \"https://cloudentity-user.authz.cloudentity.io/cloudentity-user/default\",\n  \"jti\": \"261e658f-b40a-42f5-9e98-3eb022dfccac\",\n  \"name\": \"John Doe\",\n  \"nbf\": 1631696795,\n  \"nonce\": \"c50rf23o825ulrjk38qg\",\n  \"rat\": 1631696795,\n  \"scp\": [\n    \"email\",\n    \"openid\",\n    \"profile\"\n  ],\n  \"st\": \"public\",\n  \"sub\": \"user\",\n  \"tid\": \"cloudentity-user\"\n}"

	repoMock := new(RepositoryMock)
	expected := &repository.Token{
		OidcToken: tokenValue,
	}
	repoMock.On("Store", expected).Return(nil)
	controller := New(repoMock)

	requestBody, _ := json.Marshal(&api.TokenRequest{
		OidcToken: tokenValue,
	})
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(requestBody))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)
	assert.Equal(t, 200, responseWriter.Result().StatusCode)
	repoMock.AssertExpectations(t)
}

func TestController_CreateToken_invalidJson(t *testing.T) {
	repoMock := new(RepositoryMock)
	controller := New(repoMock)

	request := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(""))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 400, status)
	repoMock.AssertNotCalled(t, "Store", mock.Anything)
}

func TestController_CreateToken_databaseError(t *testing.T) {
	tokenValue := "{\n  \"acr\": \"1\",\n  \"aid\": \"default\",\n  \"amr\": [\n    \"pwd\"\n  ],\n  \"aud\": \"default-demo\",\n  \"auth_time\": 1631696786,\n  \"email\": \"\",\n  \"email_verified\": false,\n  \"exp\": 1631700395,\n  \"iat\": 1631696795,\n  \"idp\": \"default\",\n  \"iss\": \"https://cloudentity-user.authz.cloudentity.io/cloudentity-user/default\",\n  \"jti\": \"261e658f-b40a-42f5-9e98-3eb022dfccac\",\n  \"name\": \"John Doe\",\n  \"nbf\": 1631696795,\n  \"nonce\": \"c50rf23o825ulrjk38qg\",\n  \"rat\": 1631696795,\n  \"scp\": [\n    \"email\",\n    \"openid\",\n    \"profile\"\n  ],\n  \"st\": \"public\",\n  \"sub\": \"user\",\n  \"tid\": \"cloudentity-user\"\n}"

	repoMock := new(RepositoryMock)
	repoMock.On("Store", mock.Anything).Return(errors.New("error"))

	controller := New(repoMock)
	requestBody, _ := json.Marshal(&repository.Token{
		OidcToken: tokenValue,
	})
	request := httptest.NewRequest(http.MethodPost, "/token", bytes.NewReader(requestBody))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 500, status)
}

func TestController_CreateFeedback(t *testing.T) {
	ratingComment := "any_comment"
	rating := 3

	repoMock := new(RepositoryMock)
	expected := &repository.Feedback{
		Rating:        rating,
		RatingComment: ratingComment,
		Metadata:      gormjsonb.JSONB{"first_key": "first_value", "second_key": "second_value"},
	}
	repoMock.On("Store", expected).Return(nil)
	controller := New(repoMock)

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
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)
	assert.Equal(t, 200, responseWriter.Result().StatusCode)
	repoMock.AssertExpectations(t)
}

func TestController_CreateFeedback_emptyBody(t *testing.T) {
	repoMock := new(RepositoryMock)
	controller := New(repoMock)

	request := httptest.NewRequest(http.MethodPost, "/feedback", iotest.ErrReader(errors.New("error")))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 400, status)
	repoMock.AssertNotCalled(t, "Store", mock.Anything)
}

func TestController_CreateFeedback_invalidJson(t *testing.T) {
	repoMock := new(RepositoryMock)
	controller := New(repoMock)

	request := httptest.NewRequest(http.MethodPost, "/feedback", strings.NewReader("broken"))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 400, status)
	repoMock.AssertNotCalled(t, "Store", mock.Anything)
}

func TestController_CreateFeedback_databaseError(t *testing.T) {
	repoMock := new(RepositoryMock)
	repoMock.On("Store", mock.Anything).Return(errors.New("error"))

	controller := New(repoMock)
	requestBody, _ := json.Marshal(&repository.Feedback{
		Rating:        1,
		RatingComment: "any",
	})
	request := httptest.NewRequest(http.MethodPost, "/feedback", bytes.NewReader(requestBody))
	responseWriter := httptest.NewRecorder()

	controller.GetRouter().ServeHTTP(responseWriter, request)

	status := responseWriter.Result().StatusCode
	assert.Equal(t, 500, status)
}
