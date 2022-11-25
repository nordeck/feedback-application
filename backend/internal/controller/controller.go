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
	"encoding/json"
	"feedback/internal/api"
	"feedback/internal/logger"
	"feedback/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	TOKEN_PATH    = "/token"
	FEEDBACK_PATH = "/feedback"
)

var log = logger.Instance()

type Controller struct {
	repo repository.Interface
}

func New(repo repository.Interface) *Controller {
	return &Controller{repo}
}

func (c *Controller) GetRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(TOKEN_PATH, c.createToken).Methods(http.MethodGet)
	router.HandleFunc(FEEDBACK_PATH, c.createFeedback).Methods(http.MethodPost)

	return router
}

func (c *Controller) createToken(writer http.ResponseWriter, request *http.Request) {
	authHeadVal := request.Header.Get("authorization")
	if authHeadVal == "" {
		err := "authentication header is empty!"
		http.Error(writer, err, http.StatusInternalServerError)
		log.Debug(err)
		return
	}
	dummyJwt, err := generateDummyJwt(authHeadVal)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Debug(err)
		return
	}

	err = json.NewEncoder(writer).Encode(dummyJwt)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Debug(err)
		return
	}
}

func (c *Controller) createFeedback(writer http.ResponseWriter, request *http.Request) {
	var feedback api.Feedback
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &feedback)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		log.Debug(err)
		return
	}

	err = c.repo.Store(repository.MapToFeedbackModel(feedback))

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		log.Debug(err)
		return
	}
}

func generateDummyJwt(oidcToken string) (string, error) {
	var hmacSampleSecret = []byte("SomeSampleSecret")
	trimmedHeaderValue := strings.Trim(oidcToken, "Bearer")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"oidcToken": trimmedHeaderValue,
		"nbf":       time.Date(1970, 01, 01, 0, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	return tokenString, err
}
