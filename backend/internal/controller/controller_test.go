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

func TestController_GetToken(t *testing.T) {

	repoMock := new(RepositoryMock)
	controller := New(repoMock)

	request := httptest.NewRequest(http.MethodGet, "/token", nil)
	request.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJpbmNvbWluZyIsIm5hbWUiOiJKb2huIERvZSIsImlhdCI6MTUxNjIzOTAyMn0.0DoNIGqCNa1Tc41par4bzqnQWwlgsKCIP2mgUYEHemM")

	responseWriter := httptest.NewRecorder()
	controller.GetRouter().ServeHTTP(responseWriter, request)
	assert.Equal(t, 200, responseWriter.Result().StatusCode)
	actual := responseWriter.Body.String()

	expected := "\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYmYiOjAsIm9pZGNUb2tlbiI6IiBleUpoYkdjaU9pSklVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKemRXSWlPaUpwYm1OdmJXbHVaeUlzSW01aGJXVWlPaUpLYjJodUlFUnZaU0lzSW1saGRDSTZNVFV4TmpJek9UQXlNbjAuMERvTklHcUNOYTFUYzQxcGFyNGJ6cW5RV3dsZ3NLQ0lQMm1nVVlFSGVtTSJ9.udPK9mYoV5e2MnLMwerK6j55841eimKCdf1imGVYMxg\"\n"
	assert.Equal(t, expected, actual, "not the same")
	repoMock.AssertExpectations(t)
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
